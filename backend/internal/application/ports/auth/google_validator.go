package auth

import "context"

type GoogleTokenPayload struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

type GoogleTokenValidator interface {
	Validate(ctx context.Context, idToken string) (*GoogleTokenPayload, error)
}
