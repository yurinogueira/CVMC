package mongo

import (
	"testing"
	"time"

	domainuser "cvmc/internal/domain/user"
)

func TestUserDocToDomain(t *testing.T) {
	now := time.Now().UTC()
	doc := userDoc{
		ID:              "user-123",
		Name:            "Yuri Nogueira",
		Email:           "yuri@cvmc.com",
		PasswordHash:    "hashed-password",
		GoogleID:        "google-12345",
		AvatarURL:       "https://lh3.googleusercontent.com/photo.jpg",
		EmailVerified:   true,
		EmailVerifiedAt: &now,
		MaxVehicles:     5,
		MaxFiles:        100,
		MaxStorageBytes: 100 * 1024 * 1024,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	domain := doc.toDomain()
	if domain.ID != doc.ID {
		t.Fatalf("mismatch ID: %s vs %s", domain.ID, doc.ID)
	}
	if domain.AvatarURL != doc.AvatarURL {
		t.Fatalf("mismatch AvatarURL: %s vs %s", domain.AvatarURL, doc.AvatarURL)
	}
	if domain.GoogleID != doc.GoogleID {
		t.Fatalf("mismatch GoogleID: %s vs %s", domain.GoogleID, doc.GoogleID)
	}
	if domain.Name != doc.Name || domain.Email != doc.Email {
		t.Fatalf("mismatch Name or Email")
	}
}

func TestUserDocDefaults(t *testing.T) {
	doc := userDoc{
		ID:    "user-empty",
		Email: "empty@cvmc.com",
	}

	domain := doc.toDomain()
	if domain.MaxVehicles != domainuser.DefaultMaxVehicles {
		t.Fatalf("expected default max vehicles %d, got %d", domainuser.DefaultMaxVehicles, domain.MaxVehicles)
	}
	if domain.MaxFiles != domainuser.DefaultMaxFiles {
		t.Fatalf("expected default max files %d, got %d", domainuser.DefaultMaxFiles, domain.MaxFiles)
	}
	if domain.MaxStorageBytes != domainuser.DefaultMaxStorageBytes {
		t.Fatalf("expected default max storage %d, got %d", domainuser.DefaultMaxStorageBytes, domain.MaxStorageBytes)
	}
	if domain.AvatarURL != "" {
		t.Fatalf("expected empty AvatarURL, got %q", domain.AvatarURL)
	}
}
