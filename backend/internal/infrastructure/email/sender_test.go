package email

import (
	"context"
	"errors"
	"testing"
)

func TestEmailSenderSimulation(t *testing.T) {
	service := NewService(Config{
		AppBaseURL: "https://cvmc.com.br",
		EmailFrom:  "suporte@cvmc.com.br",
	})
	ctx := context.Background()

	// 1. Send verification email
	err := service.SendVerificationEmail(ctx, "motorista@cvmc.com.br", "Yuri Nogueira", "abc123token")
	if err != nil {
		t.Fatalf("expected SendVerificationEmail to succeed in simulation mode, got %v", err)
	}

	// 2. Send password reset email
	err = service.SendPasswordResetEmail(ctx, "motorista@cvmc.com.br", "Yuri Nogueira", "reset-token-xyz")
	if err != nil {
		t.Fatalf("expected SendPasswordResetEmail to succeed in simulation mode, got %v", err)
	}
}

func TestEmailSenderRejectsInvalidAddress(t *testing.T) {
	service := NewService(Config{
		AppBaseURL: "https://cvmc.com.br",
		EmailFrom:  "suporte@cvmc.com.br",
	})
	ctx := context.Background()

	err := service.SendVerificationEmail(ctx, "invalid-email-format", "User", "token123")
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
	err := service.SendVerificationEmail(ctx, injectedRecipient, "Attacker\r\nInjected", "token\r\n123")
	if err != nil {
		// Even if parsed, the CRLF must be stripped or rejected
		t.Logf("injected recipient handled safely: %v", err)
	}

	// Direct unit tests for sanitization functions
	headerWithCRLF := "Subject with\r\nCRLF and \x00null bytes"
	sanitizedHeader := sanitizeHeader(headerWithCRLF)
	if sanitizedHeader != "Subject withCRLF and null bytes" {
		t.Fatalf("expected CRLF stripped from header, got %q", sanitizedHeader)
	}

	contentWithCRLF := "Name\r\nWith\nNewlines"
	sanitizedContent := sanitizeContent(contentWithCRLF)
	if sanitizedContent != "Name  With Newlines" {
		t.Fatalf("expected newlines replaced with spaces, got %q", sanitizedContent)
	}
}
