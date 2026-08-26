package user

import (
	"context"
	"errors"
	"testing"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

func TestUserServiceGetProfile(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	userService := NewService(users, cars, hasher)
	ctx := context.Background()

	registered, err := authService.Register(ctx, auth.RegisterInput{Name: "Roberta", Email: "roberta@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	profile, err := userService.GetProfile(ctx, registered.User.ID)
	if err != nil {
		t.Fatalf("get profile failed: %v", err)
	}
	if profile.User.Name != "Roberta" || profile.VehiclesCount != 0 || profile.MaxVehicles != 3 {
		t.Fatalf("unexpected profile output: %+v", profile)
	}
}

func TestUserServiceUpdateProfile(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	userService := NewService(users, cars, hasher)
	ctx := context.Background()

	registered, err := authService.Register(ctx, auth.RegisterInput{Name: "Roberta", Email: "roberta@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Invalid short name
	_, err = userService.UpdateProfile(ctx, registered.User.ID, "R")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for short name, got %v", err)
	}

	// Valid name
	updated, err := userService.UpdateProfile(ctx, registered.User.ID, "Roberta Silva")
	if err != nil {
		t.Fatalf("update profile failed: %v", err)
	}
	if updated.Name != "Roberta Silva" {
		t.Fatalf("unexpected updated name: %s", updated.Name)
	}
}

func TestUserServiceUpdatePassword(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	userService := NewService(users, cars, hasher)
	ctx := context.Background()

	registered, err := authService.Register(ctx, auth.RegisterInput{Name: "Roberta", Email: "roberta@example.com", Password: "initialPassword123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Invalid current password
	err = userService.UpdatePassword(ctx, registered.User.ID, "wrongPass", "newSecretPassword123")
	if !errors.Is(err, ErrInvalidCurrentPassword) {
		t.Fatalf("expected ErrInvalidCurrentPassword, got %v", err)
	}

	// Weak new password
	err = userService.UpdatePassword(ctx, registered.User.ID, "initialPassword123", "short")
	if !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword, got %v", err)
	}

	// Success update
	err = userService.UpdatePassword(ctx, registered.User.ID, "initialPassword123", "newSecretPassword123")
	if err != nil {
		t.Fatalf("update password failed: %v", err)
	}

	// Login with new password should succeed
	_, err = authService.Login(ctx, auth.LoginInput{Email: "roberta@example.com", Password: "newSecretPassword123"})
	if err != nil {
		t.Fatalf("login with new password failed: %v", err)
	}
}
