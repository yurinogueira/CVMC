package google

import (
	"context"
	"errors"
	"fmt"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	"google.golang.org/api/idtoken"
)

var (
	ErrInvalidToken        = errors.New("invalid google id token")
	ErrEmailNotVerified    = errors.New("google email is not verified")
	ErrMissingSubjectClaim = errors.New("google token missing subject claim")
	ErrMissingEmailClaim   = errors.New("google token missing email claim")
)

type Validator struct {
	clientID string
}

func NewValidator(clientID string) *Validator {
	return &Validator{
		clientID: strings.TrimSpace(clientID),
	}
}

func (v *Validator) Validate(ctx context.Context, token string) (*portauth.GoogleTokenPayload, error) {
	cleanToken := strings.TrimSpace(token)
	if cleanToken == "" {
		return nil, ErrInvalidToken
	}

	payload, err := idtoken.Validate(ctx, cleanToken, v.clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if payload.Subject == "" {
		return nil, ErrMissingSubjectClaim
	}

	email, _ := payload.Claims["email"].(string)
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, ErrMissingEmailClaim
	}

	// email_verified may come as bool or boolean string in some edge cases
	var emailVerified bool
	switch ev := payload.Claims["email_verified"].(type) {
	case bool:
		emailVerified = ev
	case string:
		emailVerified = strings.ToLower(ev) == "true"
	default:
		emailVerified = false
	}

	if !emailVerified {
		return nil, ErrEmailNotVerified
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	return &portauth.GoogleTokenPayload{
		Subject:       payload.Subject,
		Email:         email,
		EmailVerified: true,
		Name:          strings.TrimSpace(name),
		Picture:       strings.TrimSpace(picture),
	}, nil
}
