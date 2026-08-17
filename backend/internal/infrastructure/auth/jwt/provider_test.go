package jwt

import (
	"testing"

	domainuser "cvmc/internal/domain/user"
)

func TestProviderGeneratesAndParsesTokens(t *testing.T) {
	provider := NewProvider("access-secret", "refresh-secret")
	user := domainuser.User{ID: "user-1", Email: "user@example.com"}

	pair, err := provider.GeneratePair(user)
	if err != nil {
		t.Fatalf("generate pair failed: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatalf("expected tokens")
	}

	accessClaims, err := provider.ParseAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("parse access failed: %v", err)
	}
	if accessClaims.UserID != user.ID {
		t.Fatalf("unexpected user id: %s", accessClaims.UserID)
	}

	refreshClaims, err := provider.ParseRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("parse refresh failed: %v", err)
	}
	if refreshClaims.UserID != user.ID {
		t.Fatalf("unexpected refresh user id: %s", refreshClaims.UserID)
	}
}
