package notification

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	PlatformAndroid = "android"
	PlatformIOS     = "ios"
	PlatformWeb     = "web"

	MaxTitleLength = 120
	MaxBodyLength  = 500
	MaxTokenLength = 4096
)

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// DeviceToken represents a registered push notification token for a user's device.
type DeviceToken struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string    `json:"userId" bson:"userId"`
	Token     string    `json:"token" bson:"token"`
	Platform  string    `json:"platform" bson:"platform"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// PushMessage represents a sanitized notification payload with optional deep-link metadata.
type PushMessage struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

// NormalizePlatform normalizes and validates the device platform string.
func NormalizePlatform(platform string) (string, bool) {
	p := strings.ToLower(strings.TrimSpace(platform))
	if p == "" {
		return PlatformAndroid, true
	}
	switch p {
	case PlatformAndroid, PlatformIOS, PlatformWeb:
		return p, true
	default:
		return "", false
	}
}

// SanitizeText strips HTML tags, control characters, and truncates strings to maxRunes.
func SanitizeText(raw string, maxRunes int) string {
	cleaned := htmlTagRegex.ReplaceAllString(raw, "")
	var b strings.Builder
	b.Grow(len(cleaned))
	for _, r := range cleaned {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			continue
		}
		if r == '\n' || r == '\t' {
			b.WriteRune(' ')
			continue
		}
		b.WriteRune(r)
	}
	trimmed := strings.TrimSpace(strings.Join(strings.Fields(b.String()), " "))
	runes := []rune(trimmed)
	if maxRunes > 0 && len(runes) > maxRunes {
		return strings.TrimSpace(string(runes[:maxRunes]))
	}
	return trimmed
}

// SanitizeMessage returns a copy of PushMessage with sanitized title, body, and data fields.
func SanitizeMessage(msg PushMessage) PushMessage {
	sanitizedData := make(map[string]string, len(msg.Data))
	for k, v := range msg.Data {
		cleanKey := SanitizeText(k, 64)
		cleanVal := SanitizeText(v, 256)
		if cleanKey != "" {
			sanitizedData[cleanKey] = cleanVal
		}
	}
	return PushMessage{
		Title: SanitizeText(msg.Title, MaxTitleLength),
		Body:  SanitizeText(msg.Body, MaxBodyLength),
		Data:  sanitizedData,
	}
}
