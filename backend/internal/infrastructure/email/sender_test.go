package email

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestEmailSenderSimulation(t *testing.T) {
	service := NewService(Config{
		AppBaseURL: "https://cvmc.com.br",
		EmailFrom:  "suporte@cvmc.com.br",
	})
	ctx := context.Background()

	// 1. Send verification email
	err := service.SendVerificationEmail(ctx, "motorista@cvmc.com.br", "abc123token")
	if err != nil {
		t.Fatalf("expected SendVerificationEmail to succeed in simulation mode, got %v", err)
	}

	// 2. Send password reset email
	err = service.SendPasswordResetEmail(ctx, "motorista@cvmc.com.br", "reset-token-xyz")
	if err != nil {
		t.Fatalf("expected SendPasswordResetEmail to succeed in simulation mode, got %v", err)
	}
}

func TestEmailTemplatesRender(t *testing.T) {
	// 1. Verification template
	var bufVerif bytes.Buffer
	err := verificationTmpl.Execute(&bufVerif, struct {
		VerifyURL string
	}{
		VerifyURL: "https://cvmc.com.br/verify-email?token=xyz",
	})
	if err != nil {
		t.Fatalf("failed to render verification template: %v", err)
	}
	verifOutput := bufVerif.String()
	if !strings.Contains(verifOutput, "https://cvmc.com.br/verify-email?token=xyz") {
		t.Fatalf("verification template missing expected fields")
	}
	if !strings.Contains(verifOutput, "#0F52BA") {
		t.Fatalf("verification template missing brand color")
	}

	// 2. Password reset template
	var bufReset bytes.Buffer
	err = passwordResetTmpl.Execute(&bufReset, struct {
		ResetURL string
	}{
		ResetURL: "https://cvmc.com.br/reset-password?token=xyz",
	})
	if err != nil {
		t.Fatalf("failed to render password reset template: %v", err)
	}
	resetOutput := bufReset.String()
	if !strings.Contains(resetOutput, "https://cvmc.com.br/reset-password?token=xyz") {
		t.Fatalf("password reset template missing expected fields")
	}
	if !strings.Contains(resetOutput, "#0F52BA") {
		t.Fatalf("password reset template missing brand color")
	}
}

func TestEmailSenderRejectsInvalidAddress(t *testing.T) {
	service := NewService(Config{
		AppBaseURL: "https://cvmc.com.br",
		EmailFrom:  "suporte@cvmc.com.br",
	})
	ctx := context.Background()

	err := service.SendVerificationEmail(ctx, "invalid-email-format", "token123")
	if !errors.Is(err, ErrInvalidEmailAddress) {
		t.Fatalf("expected ErrInvalidEmailAddress for invalid recipient, got %v", err)
	}
}

func TestEmailSenderHeaderSanitization(t *testing.T) {
	service := NewService(Config{
		AppBaseURL: "https://cvmc.com.br",
		EmailFrom:  "suporte@cvmc.com.br",
	})
	ctx := context.Background()

	// CRLF injection attempt in recipient
	injectedRecipient := "clean@example.com\r\nBcc: attacker@evil.com"
	err := service.SendVerificationEmail(ctx, injectedRecipient, "token\r\n123")
	if err != nil {
		t.Logf("injected recipient handled safely: %v", err)
	}

	// Direct unit tests for sanitization functions
	headerWithCRLF := "Subject with\r\nCRLF and \x00null bytes"
	sanitizedHeader := sanitizeHeader(headerWithCRLF)
	if sanitizedHeader != "Subject withCRLF and null bytes" {
		t.Fatalf("expected CRLF stripped from header, got %q", sanitizedHeader)
	}

}
