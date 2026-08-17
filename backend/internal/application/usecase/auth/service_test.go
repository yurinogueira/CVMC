package auth

import (
	"context"
	"testing"

	portauth "cvmc/internal/application/ports/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestServiceRegisterLoginRefreshAndMe(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	service := NewService(users, hasher, tokens)
	ctx := context.Background()

	registered, err := service.Register(ctx, RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if registered.User.Email != "ana@example.com" {
		t.Fatalf("unexpected email: %s", registered.User.Email)
	}
	if registered.AccessToken == "" || registered.RefreshToken == "" {
		t.Fatalf("expected tokens to be generated")
	}

	logged, err := service.Login(ctx, LoginInput{Email: "ana@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if logged.User.ID != registered.User.ID {
		t.Fatalf("expected same user id, got %s and %s", logged.User.ID, registered.User.ID)
	}

	refreshed, err := service.Refresh(ctx, logged.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Fatalf("expected refreshed tokens")
	}

	me, err := service.Me(ctx, refreshed.AccessToken)
	if err != nil {
		t.Fatalf("me failed: %v", err)
	}
	if me.ID != registered.User.ID {
		t.Fatalf("unexpected user id from me: %s", me.ID)
	}
}

func TestServiceRejectsDuplicateEmail(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	service := NewService(users, hasher, tokens)
	ctx := context.Background()

	_, err := service.Register(ctx, RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	_, err = service.Register(ctx, RegisterInput{Name: "Ana 2", Email: "ana@example.com", Password: "secret123"})
	if err != ErrEmailInUse {
		t.Fatalf("expected ErrEmailInUse, got %v", err)
	}
}

var _ portauth.TokenService = (*jwtauth.Provider)(nil)
