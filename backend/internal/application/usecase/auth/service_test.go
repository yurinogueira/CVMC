package auth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	memoryuser "cvmc/internal/infrastructure/user/memory"
)

type mockEmailSender struct {
	mu                  sync.Mutex
	verificationEmails  []string
	verificationTokens  []string
	passwordResetEmails []string
	passwordResetTokens []string
}

func (m *mockEmailSender) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.verificationEmails = append(m.verificationEmails, toEmail)
	m.verificationTokens = append(m.verificationTokens, token)
	return nil
}

func (m *mockEmailSender) SendPasswordResetEmail(ctx context.Context, toEmail, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.passwordResetEmails = append(m.passwordResetEmails, toEmail)
	m.passwordResetTokens = append(m.passwordResetTokens, token)
	return nil
}

func TestServiceRegisterLoginRefreshAndMe(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	mockSender := &mockEmailSender{}
	service := NewService(users, hasher, tokens, mockSender)
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
	if len(mockSender.verificationEmails) != 1 || mockSender.verificationEmails[0] != "ana@example.com" {
		t.Fatalf("expected verification email to be sent on register")
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

func TestServiceEmailVerificationFlow(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	mockSender := &mockEmailSender{}
	service := NewService(users, hasher, tokens, mockSender)
	ctx := context.Background()

	registered, err := service.Register(ctx, RegisterInput{Name: "Carlos", Email: "carlos@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if registered.User.EmailVerified {
		t.Fatalf("expected user email to be unverified by default")
	}

	// Invalid token
	if err := service.VerifyEmail(ctx, "invalid-token"); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	// Empty token
	if err := service.VerifyEmail(ctx, ""); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for empty token, got %v", err)
	}

	// Valid verification token
	verificationToken := mockSender.verificationTokens[0]
	if err := service.VerifyEmail(ctx, verificationToken); err != nil {
		t.Fatalf("verify email failed: %v", err)
	}

	// User should now be verified
	verifiedUser, err := users.FindByID(ctx, registered.User.ID)
	if err != nil {
		t.Fatalf("user lookup failed: %v", err)
	}
	if !verifiedUser.EmailVerified || verifiedUser.EmailVerifiedAt == nil {
		t.Fatalf("expected email to be verified with timestamp")
	}

	// Using the token a second time should fail
	if err := service.VerifyEmail(ctx, verificationToken); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken when reusing verification token, got %v", err)
	}

	// Resend when already verified
	if err := service.ResendVerification(ctx, registered.User.ID); !errors.Is(err, ErrAlreadyVerified) {
		t.Fatalf("expected ErrAlreadyVerified, got %v", err)
	}
}

func TestServiceResendVerification(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	mockSender := &mockEmailSender{}
	service := NewService(users, hasher, tokens, mockSender)
	ctx := context.Background()

	registered, err := service.Register(ctx, RegisterInput{Name: "Mariana", Email: "mariana@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if err := service.ResendVerification(ctx, registered.User.ID); err != nil {
		t.Fatalf("resend verification failed: %v", err)
	}

	if len(mockSender.verificationTokens) != 2 {
		t.Fatalf("expected 2 verification tokens, got %d", len(mockSender.verificationTokens))
	}

	latestToken := mockSender.verificationTokens[1]
	if err := service.VerifyEmail(ctx, latestToken); err != nil {
		t.Fatalf("verify with resent token failed: %v", err)
	}
}

func TestServiceForgotPasswordAndResetPassword(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	mockSender := &mockEmailSender{}
	service := NewService(users, hasher, tokens, mockSender)
	ctx := context.Background()

	// Forgot password for non-existing email must return nil (anti-enumeration)
	if err := service.ForgotPassword(ctx, "nonexistent@example.com"); err != nil {
		t.Fatalf("expected nil for non-existent email, got %v", err)
	}
	if len(mockSender.passwordResetEmails) != 0 {
		t.Fatalf("expected no reset email for nonexistent user")
	}

	// Register user
	_, err := service.Register(ctx, RegisterInput{Name: "Lucas", Email: "lucas@example.com", Password: "initialPassword123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Request password reset
	if err := service.ForgotPassword(ctx, "lucas@example.com"); err != nil {
		t.Fatalf("forgot password failed: %v", err)
	}
	if len(mockSender.passwordResetEmails) != 1 || mockSender.passwordResetEmails[0] != "lucas@example.com" {
		t.Fatalf("expected reset email sent to lucas@example.com")
	}

	resetToken := mockSender.passwordResetTokens[0]

	// Reset with weak password
	if err := service.ResetPassword(ctx, resetToken, "short"); !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword, got %v", err)
	}

	// Reset with invalid token
	if err := service.ResetPassword(ctx, "invalid-token", "newValidPassword123"); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	// Successful reset
	if err := service.ResetPassword(ctx, resetToken, "newValidPassword123"); err != nil {
		t.Fatalf("reset password failed: %v", err)
	}

	// Login with old password should fail
	_, err = service.Login(ctx, LoginInput{Email: "lucas@example.com", Password: "initialPassword123"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected login with old password to fail")
	}

	// Login with new password should succeed
	_, err = service.Login(ctx, LoginInput{Email: "lucas@example.com", Password: "newValidPassword123"})
	if err != nil {
		t.Fatalf("login with new password failed: %v", err)
	}

	// Reusing reset token should fail
	if err := service.ResetPassword(ctx, resetToken, "anotherPassword123"); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken on reusing reset token")
	}
}

func TestServiceExpiredToken(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	mockSender := &mockEmailSender{}
	currentTime := time.Now().UTC()
	service := &Service{
		users:       users,
		hasher:      hasher,
		tokens:      tokens,
		emailSender: mockSender,
		now:         func() time.Time { return currentTime },
	}
	ctx := context.Background()

	_, err := service.Register(ctx, RegisterInput{Name: "Expired", Email: "expired@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	token := mockSender.verificationTokens[0]

	// Advance time by 25 hours (verification token expires in 24h)
	currentTime = currentTime.Add(25 * time.Hour)
	if err := service.VerifyEmail(ctx, token); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired for verification, got %v", err)
	}
}

func TestServiceRejectsDuplicateEmail(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	service := NewService(users, hasher, tokens, nil)
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
	service := NewService(users, hasher, tokens, nil)
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
	service := NewService(users, hasher, tokens, nil)
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
	service := NewService(users, hasher, tokens, nil)
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
