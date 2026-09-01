package config

import (
	"os"
	"strings"
	"testing"
)

func TestValidateJWTSecrets(t *testing.T) {
	validSecret := "12345678901234567890123456789012"      // 32 chars
	validRefresh := "abcdefghijklmnopqrstuvwxyz1234567890" // 36 chars

	tests := []struct {
		name          string
		secret        string
		refreshSecret string
		wantErr       bool
		errContains   string
	}{
		{
			name:          "valid secrets with 32+ chars",
			secret:        validSecret,
			refreshSecret: validRefresh,
			wantErr:       false,
		},
		{
			name:          "rejects change-me default secret",
			secret:        "change-me",
			refreshSecret: validRefresh,
			wantErr:       true,
			errContains:   "JWT_SECRET cannot use insecure default placeholder value",
		},
		{
			name:          "rejects change-me-too default refresh secret",
			secret:        validSecret,
			refreshSecret: "change-me-too",
			wantErr:       true,
			errContains:   "JWT_REFRESH_SECRET cannot use insecure default placeholder value",
		},
		{
			name:          "rejects secret shorter than 32 chars",
			secret:        "short-secret-key-12345",
			refreshSecret: validRefresh,
			wantErr:       true,
			errContains:   "JWT_SECRET must be at least 32 characters long",
		},
		{
			name:          "rejects refresh secret shorter than 32 chars",
			secret:        validSecret,
			refreshSecret: "short-refresh-key-12345",
			wantErr:       true,
			errContains:   "JWT_REFRESH_SECRET must be at least 32 characters long",
		},
		{
			name:          "rejects empty secret",
			secret:        "",
			refreshSecret: validRefresh,
			wantErr:       true,
			errContains:   "JWT_SECRET must be at least 32 characters long",
		},
		{
			name:          "rejects empty refresh secret",
			secret:        validSecret,
			refreshSecret: "",
			wantErr:       true,
			errContains:   "JWT_REFRESH_SECRET must be at least 32 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJWTSecrets(tt.secret, tt.refreshSecret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWTSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got %v", tt.errContains, err)
				}
			}
		})
	}
}

func TestConfigHTTPAddress(t *testing.T) {
	cfg := Config{Port: "9090"}
	if got := cfg.HTTPAddress(); got != ":9090" {
		t.Fatalf("HTTPAddress() = %v, want %v", got, ":9090")
	}
}

func TestGetenv(t *testing.T) {
	key := "TEST_ENV_VAR_FOR_CONFIG"
	_ = os.Unsetenv(key)

	if got := getenv(key, "default-val"); got != "default-val" {
		t.Fatalf("getenv() with unset env = %v, want %v", got, "default-val")
	}

	_ = os.Setenv(key, "custom-val")
	defer func() { _ = os.Unsetenv(key) }()

	if got := getenv(key, "default-val"); got != "custom-val" {
		t.Fatalf("getenv() with set env = %v, want %v", got, "custom-val")
	}
}

func TestParseCommaSeparated(t *testing.T) {
	items := parseCommaSeparated("http://localhost:3000, http://localhost:5173 , , https://cvmc.com")
	expected := []string{"http://localhost:3000", "http://localhost:5173", "https://cvmc.com"}

	if len(items) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(items))
	}
	for i, item := range items {
		if item != expected[i] {
			t.Errorf("item[%d] = %s, want %s", i, item, expected[i])
		}
	}
}
