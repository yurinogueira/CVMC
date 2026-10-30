package user

import (
	"net/url"
	"strings"
	"time"
)

const (
	DefaultMaxVehicles           = 3
	DefaultMaxFiles              = 100
	DefaultMaxStorageBytes int64 = 100 * 1024 * 1024 // 100 MB
)

type User struct {
	ID                         string     `json:"id"`
	Name                       string     `json:"name"`
	Email                      string     `json:"email"`
	PasswordHash               string     `json:"-"`
	GoogleID                   string     `json:"googleId,omitempty" bson:"googleId,omitempty"`
	AvatarURL                  string     `json:"avatarUrl,omitempty" bson:"avatarUrl,omitempty"`
	EmailVerified              bool       `json:"emailVerified"`
	EmailVerifiedAt            *time.Time `json:"emailVerifiedAt,omitempty"`
	EmailVerificationTokenHash string     `json:"-"`
	EmailVerificationExpiresAt *time.Time `json:"-"`
	PasswordResetTokenHash     string     `json:"-"`
	PasswordResetExpiresAt     *time.Time `json:"-"`
	MaxVehicles                int        `json:"maxVehicles"`
	FilesCount                 int        `json:"filesCount"`
	MaxFiles                   int        `json:"maxFiles"`
	StorageUsedBytes           int64      `json:"storageUsedBytes"`
	MaxStorageBytes            int64      `json:"maxStorageBytes"`
	StorageQuotaSynced         bool       `json:"-"`
	CreatedAt                  time.Time  `json:"createdAt"`
	UpdatedAt                  time.Time  `json:"updatedAt,omitempty"`
}

func (u User) HasPassword() bool {
	return u.PasswordHash != ""
}

func (u User) EffectiveMaxFiles() int {
	if u.MaxFiles <= 0 {
		return DefaultMaxFiles
	}
	return u.MaxFiles
}

func (u User) EffectiveMaxStorageBytes() int64 {
	if u.MaxStorageBytes <= 0 {
		return DefaultMaxStorageBytes
	}
	return u.MaxStorageBytes
}

func SanitizeAvatarURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return ""
	}
	if parsed.Scheme != "https" || parsed.Host == "" {
		return ""
	}
	return trimmed
}

