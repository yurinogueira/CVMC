package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/mail"
	"strings"
	"time"
	"unicode"

	portauth "cvmc/internal/application/ports/auth"
	emailport "cvmc/internal/application/ports/email"
	userport "cvmc/internal/application/ports/user"
	domainuser "cvmc/internal/domain/user"
)

const (
	MinPasswordLen = 8
	MaxPasswordLen = 72
	MinNameLen     = 2
	MaxNameLen     = 100
	MaxEmailLen    = 254
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailInUse         = errors.New("email already in use")
	ErrWeakPassword       = errors.New("weak password")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAlreadyVerified    = errors.New("email already verified")
)

type Service struct {
	users       userport.Repository
	hasher      portauth.PasswordHasher
	tokens      portauth.TokenService
	emailSender emailport.Sender
	now         func() time.Time
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	User         domainuser.User
	AccessToken  string
	RefreshToken string
}

func NewService(users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService, emailSender emailport.Sender) *Service {
	return &Service{
		users:       users,
		hasher:      hasher,
		tokens:      tokens,
		emailSender: emailSender,
		now:         time.Now,
	}
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func generateCryptoToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > MaxEmailLen {
		return false
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	domainParts := strings.Split(parts[1], ".")
	if len(domainParts) < 2 {
		return false
	}
	for _, dp := range domainParts {
		if dp == "" {
			return false
		}
	}
	return true
}

func isValidName(name string) bool {
	if len(name) < MinNameLen || len(name) > MaxNameLen {
		return false
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthOutput, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Email == "" || input.Name == "" || input.Password == "" {
		return AuthOutput{}, ErrInvalidInput
	}
	if !isValidName(input.Name) {
		return AuthOutput{}, ErrInvalidInput
	}
	if !isValidEmail(input.Email) {
		return AuthOutput{}, ErrInvalidInput
	}
	if len(input.Password) < MinPasswordLen || len(input.Password) > MaxPasswordLen {
		return AuthOutput{}, ErrWeakPassword
	}
	if _, err := s.users.FindByEmail(ctx, input.Email); err == nil {
		return AuthOutput{}, ErrEmailInUse
	}
	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return AuthOutput{}, err
	}

	verificationToken, err := generateCryptoToken()
	if err != nil {
		return AuthOutput{}, err
	}
	verificationTokenHash := hashToken(verificationToken)
	verificationExpiry := s.now().UTC().Add(24 * time.Hour)

	created, err := s.users.Create(ctx, domainuser.User{
		Name:                       input.Name,
		Email:                      input.Email,
		PasswordHash:               hash,
		EmailVerified:              false,
		EmailVerificationTokenHash: verificationTokenHash,
		EmailVerificationExpiresAt: &verificationExpiry,
		MaxVehicles:                3,
		CreatedAt:                  s.now().UTC(),
	})
	if err != nil {
		return AuthOutput{}, err
	}

	if s.emailSender != nil {
		if err := s.emailSender.SendVerificationEmail(ctx, created.Email, verificationToken); err != nil {
			log.Printf("[AUTH] Warning: failed to send verification email to %s: %v", created.Email, err)
		}
	}

	pair, err := s.tokens.GeneratePair(created)
	if err != nil {
		return AuthOutput{}, err
	}
	return AuthOutput{User: created, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthOutput, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if input.Email == "" || input.Password == "" {
		return AuthOutput{}, ErrInvalidCredentials
	}
	if len(input.Email) > MaxEmailLen || len(input.Password) > MaxPasswordLen {
		return AuthOutput{}, ErrInvalidCredentials
	}
	user, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		return AuthOutput{}, ErrInvalidCredentials
	}
	if err := s.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return AuthOutput{}, ErrInvalidCredentials
	}
	pair, err := s.tokens.GeneratePair(user)
	if err != nil {
		return AuthOutput{}, err
	}
	return AuthOutput{User: user, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (AuthOutput, error) {
	claims, err := s.tokens.ParseRefreshToken(strings.TrimSpace(refreshToken))
	if err != nil {
		return AuthOutput{}, ErrInvalidToken
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return AuthOutput{}, ErrUserNotFound
	}
	pair, err := s.tokens.GeneratePair(user)
	if err != nil {
		return AuthOutput{}, err
	}
	return AuthOutput{User: user, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Service) Me(ctx context.Context, accessToken string) (domainuser.User, error) {
	claims, err := s.tokens.ParseAccessToken(strings.TrimSpace(accessToken))
	if err != nil {
		return domainuser.User{}, ErrInvalidToken
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return domainuser.User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	cleanToken := strings.TrimSpace(token)
	if cleanToken == "" {
		return ErrInvalidToken
	}

	tokenHash := hashToken(cleanToken)
	user, err := s.users.FindByEmailVerificationTokenHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidToken
	}

	now := s.now().UTC()
	if user.EmailVerificationExpiresAt == nil || user.EmailVerificationExpiresAt.Before(now) {
		return ErrTokenExpired
	}

	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.EmailVerificationTokenHash = ""
	user.EmailVerificationExpiresAt = nil
	user.UpdatedAt = now

	_, err = s.users.Update(ctx, user)
	return err
}

func (s *Service) ResendVerification(ctx context.Context, userID string) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.EmailVerified {
		return ErrAlreadyVerified
	}

	verificationToken, err := generateCryptoToken()
	if err != nil {
		return err
	}
	now := s.now().UTC()
	expiry := now.Add(24 * time.Hour)

	user.EmailVerificationTokenHash = hashToken(verificationToken)
	user.EmailVerificationExpiresAt = &expiry
	user.UpdatedAt = now

	if _, err := s.users.Update(ctx, user); err != nil {
		return err
	}

	if s.emailSender != nil {
		return s.emailSender.SendVerificationEmail(ctx, user.Email, verificationToken)
	}
	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	if !isValidEmail(cleanEmail) {
		return ErrInvalidInput
	}

	user, err := s.users.FindByEmail(ctx, cleanEmail)
	if err != nil {
		// Anti-enumeration: Return nil as success even if user not found
		return nil
	}

	resetToken, err := generateCryptoToken()
	if err != nil {
		return err
	}
	now := s.now().UTC()
	expiry := now.Add(30 * time.Minute)

	user.PasswordResetTokenHash = hashToken(resetToken)
	user.PasswordResetExpiresAt = &expiry
	user.UpdatedAt = now

	if _, err := s.users.Update(ctx, user); err != nil {
		return err
	}

	if s.emailSender != nil {
		if err := s.emailSender.SendPasswordResetEmail(ctx, user.Email, resetToken); err != nil {
			log.Printf("[AUTH] Warning: failed to send password reset email to %s: %v", user.Email, err)
		}
	}
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	cleanToken := strings.TrimSpace(token)
	if cleanToken == "" {
		return ErrInvalidToken
	}
	if len(newPassword) < MinPasswordLen || len(newPassword) > MaxPasswordLen {
		return ErrWeakPassword
	}

	tokenHash := hashToken(cleanToken)
	user, err := s.users.FindByPasswordResetTokenHash(ctx, tokenHash)
	if err != nil {
		return ErrInvalidToken
	}

	now := s.now().UTC()
	if user.PasswordResetExpiresAt == nil || user.PasswordResetExpiresAt.Before(now) {
		return ErrTokenExpired
	}

	newHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	user.PasswordResetTokenHash = ""
	user.PasswordResetExpiresAt = nil
	user.UpdatedAt = now

	_, err = s.users.Update(ctx, user)
	return err
}
