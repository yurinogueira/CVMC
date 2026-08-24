package auth

import (
	"context"
	"errors"
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

func TestServiceRejectsWeakPassword(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	service := NewService(users, hasher, tokens)
	ctx := context.Background()

	// Less than 8 chars
	_, err := service.Register(ctx, RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "short"})
	if !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword for short password, got %v", err)
	}

	// Greater than 72 chars (bcrypt DoS limit)
	longPassword := "1234567890123456789012345678901234567890123456789012345678901234567890123" // 73 chars
	_, err = service.Register(ctx, RegisterInput{Name: "Ana", Email: "ana@example.com", Password: longPassword})
	if !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword for password > 72 chars, got %v", err)
	}

	// Exact boundaries: 8 chars and 72 chars should succeed
	exact8 := "12345678"
	exact72 := "123456789012345678901234567890123456789012345678901234567890123456789012"
	_, err = service.Register(ctx, RegisterInput{Name: "User Eight", Email: "user8@example.com", Password: exact8})
	if err != nil {
		t.Fatalf("expected success for 8 char password, got %v", err)
	}
	_, err = service.Register(ctx, RegisterInput{Name: "User 72", Email: "user72@example.com", Password: exact72})
	if err != nil {
		t.Fatalf("expected success for 72 char password, got %v", err)
	}
}

func TestServiceRejectsInvalidNameAndEmail(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	service := NewService(users, hasher, tokens)
	ctx := context.Background()

	// Short name (< 2 chars)
	_, err := service.Register(ctx, RegisterInput{Name: "A", Email: "valid@example.com", Password: "password123"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for 1-char name, got %v", err)
	}

	// Long name (> 100 chars)
	longName := string(make([]byte, 101))
	_, err = service.Register(ctx, RegisterInput{Name: longName, Email: "valid@example.com", Password: "password123"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for name > 100 chars, got %v", err)
	}

	// Invalid email formats
	invalidEmails := []string{
		"plainaddress",
		"@missingusername.com",
		"username@.com",
		"username@domain",
		"user name@example.com",
	}
	for _, email := range invalidEmails {
		_, err = service.Register(ctx, RegisterInput{Name: "Valid Name", Email: email, Password: "password123"})
		if !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("expected ErrInvalidInput for email '%s', got %v", email, err)
		}
	}
}

func TestServiceLoginRejectsOversizedInputs(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	service := NewService(users, hasher, tokens)
	ctx := context.Background()

	// Oversized password on login (> 72 chars)
	oversizedPass := string(make([]byte, 500))
	_, err := service.Login(ctx, LoginInput{Email: "test@example.com", Password: oversizedPass})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for oversized password on login, got %v", err)
	}

	// Oversized email on login (> 254 chars)
	oversizedEmail := string(make([]byte, 300))
	_, err = service.Login(ctx, LoginInput{Email: oversizedEmail, Password: "password123"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for oversized email on login, got %v", err)
	}
}

var _ portauth.TokenService = (*jwtauth.Provider)(nil)
