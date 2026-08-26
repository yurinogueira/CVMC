package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
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

func TestEmailTemplatesRender(t *testing.T) {
	currentYear := time.Now().Year()
	expectedYearStr := fmt.Sprintf("© %d CVMC", currentYear)

	// 1. Verification template
	var bufVerif bytes.Buffer
	err := verificationTmpl.Execute(&bufVerif, struct {
		Name        string
		VerifyURL   string
		CurrentYear int
	}{
		Name:        "Carlos Silva",
		VerifyURL:   "https://cvmc.com.br/verify-email?token=xyz",
		CurrentYear: currentYear,
	})
	if err != nil {
		t.Fatalf("failed to render verification template: %v", err)
	}
	verifStr := bufVerif.String()
	if !strings.Contains(verifStr, "Carlos Silva") || !strings.Contains(verifStr, "https://cvmc.com.br/verify-email?token=xyz") {
		t.Fatalf("verification template missing expected fields")
	}
	if !strings.Contains(verifStr, expectedYearStr) {
		t.Fatalf("verification template missing dynamic current year, got: %s", verifStr)
	}
	if !strings.Contains(verifStr, "#0F52BA") {
		t.Fatalf("verification template missing brand color")
	}

	// 2. Password reset template
	var bufReset bytes.Buffer
	err = passwordResetTmpl.Execute(&bufReset, struct {
		Name        string
		ResetURL    string
		CurrentYear int
	}{
		Name:        "Carlos Silva",
		ResetURL:    "https://cvmc.com.br/reset-password?token=xyz",
		CurrentYear: currentYear,
	})
	if err != nil {
		t.Fatalf("failed to render password reset template: %v", err)
	}
	resetStr := bufReset.String()
	if !strings.Contains(resetStr, "Carlos Silva") || !strings.Contains(resetStr, "https://cvmc.com.br/reset-password?token=xyz") {
		t.Fatalf("password reset template missing expected fields")
	}
	if !strings.Contains(resetStr, expectedYearStr) {
		t.Fatalf("password reset template missing dynamic current year, got: %s", resetStr)
	}
	if !strings.Contains(resetStr, "#0F52BA") {
		t.Fatalf("password reset template missing brand color")
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
