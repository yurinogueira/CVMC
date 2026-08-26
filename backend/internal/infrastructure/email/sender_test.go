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

	err := service.SendVerificationEmail(ctx, "motorista@cvmc.com.br", "abc123token")
	if err != nil {
		t.Fatalf("expected SendVerificationEmail to succeed in simulation mode, got %v", err)
	}

	err = service.SendPasswordResetEmail(ctx, "motorista@cvmc.com.br", "reset-token-xyz")
	if err != nil {
		t.Fatalf("expected SendPasswordResetEmail to succeed in simulation mode, got %v", err)
	}
}

func TestEmailTemplatesRender(t *testing.T) {
	var bufVerif bytes.Buffer
	err := verificationTmpl.Execute(&bufVerif, struct{ VerifyURL string }{
		VerifyURL: "https://cvmc.com.br/verify-email?token=xyz",
	})
	if err != nil {
		t.Fatalf("failed to render verification template: %v", err)
	}
	if !strings.Contains(bufVerif.String(), "https://cvmc.com.br/verify-email?token=xyz") {
		t.Fatalf("verification template missing expected URL")
	}
	if !strings.Contains(bufVerif.String(), "#0F52BA") {
		t.Fatalf("verification template missing brand color")
	}

	var bufReset bytes.Buffer
	err = passwordResetTmpl.Execute(&bufReset, struct{ ResetURL string }{
		ResetURL: "https://cvmc.com.br/reset-password?token=xyz",
	})
	if err != nil {
		t.Fatalf("failed to render password reset template: %v", err)
	}
	if !strings.Contains(bufReset.String(), "https://cvmc.com.br/reset-password?token=xyz") {
		t.Fatalf("password reset template missing expected URL")
	}
	if !strings.Contains(bufReset.String(), "#0F52BA") {
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

func TestValidateEmailRejectsCRLF(t *testing.T) {
	cases := []string{
		"user@example.com\r\nBcc: attacker@evil.com",
		"user@example.com\nBcc: attacker@evil.com",
		"user\r@example.com",
		"not-an-email",
		"",
		"user@",
		"@domain.com",
	}
	for _, addr := range cases {
		if _, err := validateEmail(addr); err == nil {
			t.Errorf("expected validateEmail to reject %q, but it passed", addr)
		}
	}
}

func TestValidateEmailAcceptsValid(t *testing.T) {
	cases := []string{
		"user@example.com",
		"first.last@domain.co",
		"test+tag@sub.domain.com",
		"a@b.cd",
	}
	for _, addr := range cases {
		if _, err := validateEmail(addr); err != nil {
			t.Errorf("expected validateEmail to accept %q, got %v", addr, err)
		}
	}
}
