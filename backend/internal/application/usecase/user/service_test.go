package user

import (
	"context"
	"errors"
	"testing"

	portauth "cvmc/internal/application/ports/auth"
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

type mockGoogleValidator struct {
	payload *portauth.GoogleTokenPayload
	err     error
}

func (m *mockGoogleValidator) Validate(ctx context.Context, idToken string) (*portauth.GoogleTokenPayload, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.payload, nil
}

func TestUserServiceAuthProvidersAndLinkUnlinkGoogle(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret-32-chars-long-security", "refresh-secret-32-chars-long-security")
	mockValidator := &mockGoogleValidator{
		payload: &portauth.GoogleTokenPayload{
			Subject:       "google-id-alice",
			Email:         "alice@example.com",
			EmailVerified: true,
			Name:          "Alice",
		},
	}
	authService := auth.NewService(users, hasher, tokens, nil)
	userService := NewService(users, cars, hasher, mockValidator)
	ctx := context.Background()

	reg, err := authService.Register(ctx, auth.RegisterInput{Name: "Alice", Email: "alice@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("failed to register alice: %v", err)
	}

	// 1. Initial providers check
	providers, err := userService.GetAuthProviders(ctx, reg.User.ID)
	if err != nil {
		t.Fatalf("failed to get providers: %v", err)
	}
	if !providers.HasPassword || providers.GoogleLinked {
		t.Fatalf("unexpected providers: %+v", providers)
	}

	// 2. Link Google with mismatched email fails
	mockValidator.payload.Email = "other@example.com"
	err = userService.LinkGoogleAccount(ctx, reg.User.ID, "token-1")
	if !errors.Is(err, ErrGoogleEmailMismatch) {
		t.Fatalf("expected ErrGoogleEmailMismatch, got %v", err)
	}

	// 3. Link Google with matching email succeeds
	mockValidator.payload.Email = "alice@example.com"
	err = userService.LinkGoogleAccount(ctx, reg.User.ID, "token-1")
	if err != nil {
		t.Fatalf("expected successful link, got %v", err)
	}

	providers, err = userService.GetAuthProviders(ctx, reg.User.ID)
	if err != nil {
		t.Fatalf("failed to get providers: %v", err)
	}
	if !providers.GoogleLinked || providers.GoogleEmail != "alice@example.com" {
		t.Fatalf("expected google to be linked: %+v", providers)
	}

	// 4. Another user tries to link the same Google ID fails
	regBob, _ := authService.Register(ctx, auth.RegisterInput{Name: "Bob", Email: "bob@example.com", Password: "password123"})
	mockValidator.payload.Email = "bob@example.com"
	err = userService.LinkGoogleAccount(ctx, regBob.User.ID, "token-2")
	if !errors.Is(err, ErrGoogleAlreadyLinked) {
		t.Fatalf("expected ErrGoogleAlreadyLinked, got %v", err)
	}

	// 5. Unlink Google succeeds for Alice (has password)
	err = userService.UnlinkGoogleAccount(ctx, reg.User.ID)
	if err != nil {
		t.Fatalf("expected successful unlink, got %v", err)
	}
	providers, _ = userService.GetAuthProviders(ctx, reg.User.ID)
	if providers.GoogleLinked {
		t.Fatalf("expected google to be unlinked")
	}
}

func TestUserServiceGoogleOnlyUser_LockoutPreventionAndSetPassword(t *testing.T) {
	users := memoryuser.NewRepository()
	cars := carrepo.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret-32-chars-long-security", "refresh-secret-32-chars-long-security")
	mockValidator := &mockGoogleValidator{
		payload: &portauth.GoogleTokenPayload{
			Subject:       "google-id-carol",
			Email:         "carol@example.com",
			EmailVerified: true,
			Name:          "Carol",
		},
	}
	authService := auth.NewService(users, hasher, tokens, nil, mockValidator)
	userService := NewService(users, cars, hasher, mockValidator)
	ctx := context.Background()

	// Register Carol via Google login (no password)
	res, err := authService.LoginWithGoogle(ctx, "carol-token")
	if err != nil {
		t.Fatalf("failed google login: %v", err)
	}

	// Check providers: hasPassword should be false
	providers, err := userService.GetAuthProviders(ctx, res.User.ID)
	if err != nil {
		t.Fatalf("failed to get providers: %v", err)
	}
	if providers.HasPassword || !providers.GoogleLinked {
		t.Fatalf("expected no password and google linked: %+v", providers)
	}

	// Unlink Google should fail to prevent lockout
	err = userService.UnlinkGoogleAccount(ctx, res.User.ID)
	if !errors.Is(err, ErrLockoutPrevented) {
		t.Fatalf("expected ErrLockoutPrevented, got %v", err)
	}

	// Set initial password
	err = userService.SetPassword(ctx, res.User.ID, "carolNewPass123")
	if err != nil {
		t.Fatalf("failed to set password: %v", err)
	}

	// Setting password again fails
	err = userService.SetPassword(ctx, res.User.ID, "anotherPass123")
	if !errors.Is(err, ErrPasswordAlreadySet) {
		t.Fatalf("expected ErrPasswordAlreadySet, got %v", err)
	}

	// Now Unlink Google should succeed
	err = userService.UnlinkGoogleAccount(ctx, res.User.ID)
	if err != nil {
		t.Fatalf("expected unlink to succeed now that password is set, got %v", err)
	}
}
