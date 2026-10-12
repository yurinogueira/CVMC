package google

import (
	"context"
	"errors"
	"testing"
)

func TestValidator_UnconfiguredClientID(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name     string
		clientID string
	}{
		{name: "empty string", clientID: ""},
		{name: "only spaces", clientID: "   "},
		{name: "tabs and spaces", clientID: "\t  \n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := NewValidator(tc.clientID)
			_, err := v.Validate(ctx, "sample.google.id.token")
			if !errors.Is(err, ErrGoogleAuthNotConfigured) {
				t.Fatalf("expected ErrGoogleAuthNotConfigured, got %v", err)
			}
		})
	}
}

func TestValidator_EmptyOrInvalidToken(t *testing.T) {
	ctx := context.Background()
	v := NewValidator("1234567890-sample.apps.googleusercontent.com")

	t.Run("empty token", func(t *testing.T) {
		_, err := v.Validate(ctx, "")
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("whitespace token", func(t *testing.T) {
		_, err := v.Validate(ctx, "   \t\n")
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("malformed token format", func(t *testing.T) {
		_, err := v.Validate(ctx, "not-a-valid-jwt")
		if err == nil {
			t.Fatal("expected error for malformed token, got nil")
		}
		if !errors.Is(err, ErrInvalidToken) {
			t.Fatalf("expected error wrapping ErrInvalidToken, got %v", err)
		}
	})
}
