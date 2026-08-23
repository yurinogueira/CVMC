package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	userport "cvmc/internal/application/ports/user"
	domainuser "cvmc/internal/domain/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailInUse         = errors.New("email already in use")
	ErrWeakPassword       = errors.New("weak password")
)

type Service struct {
	users  userport.Repository
	hasher portauth.PasswordHasher
	tokens portauth.TokenService
	now    func() time.Time
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

func NewService(users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService) *Service {
	return &Service{users: users, hasher: hasher, tokens: tokens, now: time.Now}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthOutput, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Email == "" || input.Name == "" || input.Password == "" {
		return AuthOutput{}, ErrInvalidCredentials
	}
	if len(input.Password) < 8 {
		return AuthOutput{}, ErrWeakPassword
	}
	if _, err := s.users.FindByEmail(ctx, input.Email); err == nil {
		return AuthOutput{}, ErrEmailInUse
	}
	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return AuthOutput{}, err
	}
	created, err := s.users.Create(ctx, domainuser.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		CreatedAt:    s.now().UTC(),
	})
	if err != nil {
		return AuthOutput{}, err
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
