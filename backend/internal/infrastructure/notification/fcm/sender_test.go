package fcm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	notificationport "cvmc/internal/application/ports/notification"
	domainnotification "cvmc/internal/domain/notification"

	"golang.org/x/oauth2"
)

type staticTokenSource struct {
	token string
}

func (s staticTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: s.token}, nil
}

func TestSender_NoopWhenNoCredentials(t *testing.T) {
	sender, err := NewSender(context.Background(), Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = sender.Send(context.Background(), "token-123456789", domainnotification.PushMessage{
		Title: "Bem-vindo ao CVMC",
		Body:  "Notificações ativadas!",
	})
	if err != nil {
		t.Fatalf("expected nil error in noop mode, got: %v", err)
	}
}

func TestSender_SendSuccess(t *testing.T) {
	var received fcmSendRequest
	var authHeader string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/cvmc-test/messages:send" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		authHeader = r.Header.Get("Authorization")
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"projects/cvmc-test/messages/123"}`))
	}))
	defer srv.Close()

	sender, err := NewSender(context.Background(), Config{
		ProjectID:   "cvmc-test",
		BaseURL:     srv.URL,
		HTTPClient:  srv.Client(),
		TokenSource: staticTokenSource{token: "fake-oauth-token"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = sender.Send(context.Background(), "device-token-abc", domainnotification.PushMessage{
		Title: "<b>Revisão Preventiva</b>",
		Body:  "Seu Civic está próximo da revisão.",
		Data: map[string]string{
			"route": "/vehicles/car-1",
			"carId": "car-1",
		},
	})
	if err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
	if authHeader != "Bearer fake-oauth-token" {
		t.Errorf("expected Bearer fake-oauth-token, got %q", authHeader)
	}
	if received.Message.Token != "device-token-abc" {
		t.Errorf("expected token device-token-abc, got %q", received.Message.Token)
	}
	if received.Message.Notification.Title != "Revisão Preventiva" {
		t.Errorf("expected sanitized title 'Revisão Preventiva', got %q", received.Message.Notification.Title)
	}
	if received.Message.Data["route"] != "/vehicles/car-1" {
		t.Errorf("expected deep link route '/vehicles/car-1', got %q", received.Message.Data["route"])
	}
	if received.Message.Android.Notification.ChannelID != defaultAndroidChannel {
		t.Errorf("expected channel %q, got %q", defaultAndroidChannel, received.Message.Android.Notification.ChannelID)
	}
}

func TestSender_UnregisteredTokenDetection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{
			"error": {
				"code": 404,
				"message": "Requested entity was not found.",
				"status": "NOT_FOUND",
				"details": [
					{
						"@type": "type.googleapis.com/google.firebase.fcm.v1.FcmError",
						"errorCode": "UNREGISTERED"
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	sender, err := NewSender(context.Background(), Config{
		ProjectID:   "cvmc-test",
		BaseURL:     srv.URL,
		HTTPClient:  srv.Client(),
		TokenSource: staticTokenSource{token: "fake-oauth-token"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = sender.Send(context.Background(), "expired-device-token", domainnotification.PushMessage{
		Title: "Teste",
		Body:  "Mensagem",
	})
	if !errors.Is(err, notificationport.ErrUnregisteredToken) {
		t.Fatalf("expected ErrUnregisteredToken, got: %v", err)
	}
}
