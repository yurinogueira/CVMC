package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	notificationport "cvmc/internal/application/ports/notification"
	domainnotification "cvmc/internal/domain/notification"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	fcmMessagingScope     = "https://www.googleapis.com/auth/firebase.messaging"
	defaultFCMBaseURL     = "https://fcm.googleapis.com/v1"
	defaultAndroidChannel = "cvmc_alerts"
)

type Config struct {
	CredentialsJSON string
	CredentialsFile string
	BaseURL         string
	HTTPClient      *http.Client
	TokenSource     oauth2.TokenSource
	ProjectID       string
}

type Sender struct {
	projectID   string
	baseURL     string
	httpClient  *http.Client
	tokenSource oauth2.TokenSource
	noop        bool
}

type serviceAccountInfo struct {
	ProjectID string `json:"project_id"`
}

func NewSender(ctx context.Context, cfg Config) (*Sender, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultFCMBaseURL
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	// Allow explicit injection for tests
	if cfg.TokenSource != nil && strings.TrimSpace(cfg.ProjectID) != "" {
		return &Sender{
			projectID:   strings.TrimSpace(cfg.ProjectID),
			baseURL:     baseURL,
			httpClient:  httpClient,
			tokenSource: cfg.TokenSource,
		}, nil
	}

	var credBytes []byte
	if strings.TrimSpace(cfg.CredentialsJSON) != "" {
		credBytes = []byte(strings.TrimSpace(cfg.CredentialsJSON))
	} else if strings.TrimSpace(cfg.CredentialsFile) != "" {
		b, err := os.ReadFile(strings.TrimSpace(cfg.CredentialsFile))
		if err != nil {
			return nil, fmt.Errorf("failed to read firebase credentials file: %w", err)
		}
		credBytes = b
	}

	if len(credBytes) == 0 {
		return &Sender{
			baseURL:    baseURL,
			httpClient: httpClient,
			noop:       true,
		}, nil
	}

	var sa serviceAccountInfo
	if err := json.Unmarshal(credBytes, &sa); err != nil {
		return nil, fmt.Errorf("invalid firebase credentials json: %w", err)
	}
	projectID := strings.TrimSpace(cfg.ProjectID)
	if projectID == "" {
		projectID = strings.TrimSpace(sa.ProjectID)
	}
	if projectID == "" {
		return nil, errors.New("firebase credentials missing project_id")
	}

	creds, err := google.CredentialsFromJSON(ctx, credBytes, fcmMessagingScope)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize google credentials for fcm: %w", err)
	}

	return &Sender{
		projectID:   projectID,
		baseURL:     baseURL,
		httpClient:  httpClient,
		tokenSource: creds.TokenSource,
	}, nil
}

type fcmSendRequest struct {
	Message fcmMessage `json:"message"`
}

type fcmMessage struct {
	Token        string            `json:"token"`
	Notification *fcmNotification  `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Android      *fcmAndroidConfig `json:"android,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type fcmAndroidConfig struct {
	Priority     string                  `json:"priority,omitempty"`
	Notification *fcmAndroidNotification `json:"notification,omitempty"`
}

type fcmAndroidNotification struct {
	ChannelID string `json:"channel_id,omitempty"`
	Sound     string `json:"sound,omitempty"`
}

type fcmErrorEnvelope struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
		Details []struct {
			Type      string `json:"@type"`
			ErrorCode string `json:"errorCode"`
		} `json:"details"`
	} `json:"error"`
}

func (s *Sender) Send(ctx context.Context, token string, msg domainnotification.PushMessage) error {
	cleanToken := strings.TrimSpace(token)
	if cleanToken == "" {
		return errors.New("empty device token")
	}

	sanitized := domainnotification.SanitizeMessage(msg)

	if s.noop {
		log.Printf("[FCM NO-OP] Push notification to token %s: title=%q body=%q data=%v", maskToken(cleanToken), sanitized.Title, sanitized.Body, sanitized.Data)
		return nil
	}

	oauthToken, err := s.tokenSource.Token()
	if err != nil {
		return fmt.Errorf("fcm oauth token error: %w", err)
	}

	payload := fcmSendRequest{
		Message: fcmMessage{
			Token: cleanToken,
			Notification: &fcmNotification{
				Title: sanitized.Title,
				Body:  sanitized.Body,
			},
			Data: sanitized.Data,
			Android: &fcmAndroidConfig{
				Priority: "HIGH",
				Notification: &fcmAndroidNotification{
					ChannelID: defaultAndroidChannel,
					Sound:     "default",
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal fcm payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/projects/%s/messages:send", s.baseURL, s.projectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create fcm request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+oauthToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fcm http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	if isUnregisteredError(resp.StatusCode, respBody) {
		return notificationport.ErrUnregisteredToken
	}

	return fmt.Errorf("fcm send failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
}

func isUnregisteredError(statusCode int, body []byte) bool {
	if statusCode == http.StatusNotFound || statusCode == http.StatusGone {
		return true
	}
	var env fcmErrorEnvelope
	if err := json.Unmarshal(body, &env); err == nil {
		if strings.EqualFold(env.Error.Status, "NOT_FOUND") || strings.EqualFold(env.Error.Status, "UNREGISTERED") {
			return true
		}
		for _, detail := range env.Error.Details {
			if strings.EqualFold(detail.ErrorCode, "UNREGISTERED") {
				return true
			}
		}
	}
	return strings.Contains(strings.ToUpper(string(body)), "UNREGISTERED")
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
