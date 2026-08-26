package user

import "time"

type User struct {
	ID                         string     `json:"id"`
	Name                       string     `json:"name"`
	Email                      string     `json:"email"`
	PasswordHash               string     `json:"-"`
	EmailVerified              bool       `json:"emailVerified"`
	EmailVerifiedAt            *time.Time `json:"emailVerifiedAt,omitempty"`
	EmailVerificationTokenHash string     `json:"-"`
	EmailVerificationExpiresAt *time.Time `json:"-"`
	PasswordResetTokenHash     string     `json:"-"`
	PasswordResetExpiresAt     *time.Time `json:"-"`
	MaxVehicles                int        `json:"maxVehicles"`
	CreatedAt                  time.Time  `json:"createdAt"`
	UpdatedAt                  time.Time  `json:"updatedAt,omitempty"`
}
