package email

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	emailport "cvmc/internal/application/ports/email"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	ErrInvalidEmailAddress = errors.New("invalid email address")

	verificationTmpl  = htmltemplate.Must(htmltemplate.ParseFS(templateFS, "templates/verification.html"))
	passwordResetTmpl = htmltemplate.Must(htmltemplate.ParseFS(templateFS, "templates/password_reset.html"))
)

type Config struct {
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
	EmailFrom  string
	AppBaseURL string
}

type Service struct {
	cfg Config
}

func NewService(cfg Config) emailport.Sender {
	if cfg.AppBaseURL == "" {
		cfg.AppBaseURL = "http://localhost:5173"
	}
	if cfg.EmailFrom == "" {
		cfg.EmailFrom = "no-reply@cvmc.com.br"
	}
	return &Service{cfg: cfg}
}

func (s *Service) SendVerificationEmail(ctx context.Context, toEmail, toName, token string) error {
	_ = ctx

	cleanName := strings.TrimSpace(toName)
	if cleanName == "" {
		cleanName = "Motorista"
	}

	safeToken := url.QueryEscape(strings.TrimSpace(token))
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", strings.TrimRight(s.cfg.AppBaseURL, "/"), safeToken)

	subject := "CVMC - Confirmação de E-mail"

	var htmlBuf bytes.Buffer
	data := struct {
		Name        string
		VerifyURL   string
		CurrentYear int
	}{
		Name:        cleanName,
		VerifyURL:   verifyURL,
		CurrentYear: time.Now().Year(),
	}
	if err := verificationTmpl.Execute(&htmlBuf, data); err != nil {
		return fmt.Errorf("failed to render verification email template: %w", err)
	}

	plainBody := fmt.Sprintf(
		"Olá, %s!\n\nObrigado por se cadastrar no CVMC (Como Vai Meu Carro).\n\nPara validar seu e-mail e liberar o cadastro de veículos, acesse o link abaixo:\n%s\n\nEste link é válido por 24 horas.\n\nSe você não criou uma conta no CVMC, ignore este e-mail.",
		cleanName,
		verifyURL,
	)

	return s.send(toEmail, subject, plainBody, htmlBuf.String())
}

func (s *Service) SendPasswordResetEmail(ctx context.Context, toEmail, toName, token string) error {
	_ = ctx

	cleanName := strings.TrimSpace(toName)
	if cleanName == "" {
		cleanName = "Motorista"
	}

	safeToken := url.QueryEscape(strings.TrimSpace(token))
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimRight(s.cfg.AppBaseURL, "/"), safeToken)

	subject := "CVMC - Recuperação de Senha"

	var htmlBuf bytes.Buffer
	data := struct {
		Name        string
		ResetURL    string
		CurrentYear int
	}{
		Name:        cleanName,
		ResetURL:    resetURL,
		CurrentYear: time.Now().Year(),
	}
	if err := passwordResetTmpl.Execute(&htmlBuf, data); err != nil {
		return fmt.Errorf("failed to render password reset email template: %w", err)
	}

	plainBody := fmt.Sprintf(
		"Olá, %s!\n\nRecebemos uma solicitação para redefinir a senha da sua conta no CVMC.\n\nPara criar uma nova senha, acesse o link abaixo:\n%s\n\nEste link é válido por 30 minutos.\n\nSe você não solicitou a redefinição de senha, ignore este e-mail com segurança.",
		cleanName,
		resetURL,
	)

	return s.send(toEmail, subject, plainBody, htmlBuf.String())
}

func (s *Service) send(toEmail, subject, plainBody, htmlBody string) error {
	parsedTo, err := mail.ParseAddress(strings.TrimSpace(toEmail))
	if err != nil {
		return fmt.Errorf("%w: recipient %q", ErrInvalidEmailAddress, toEmail)
	}

	parsedFrom, err := mail.ParseAddress(strings.TrimSpace(s.cfg.EmailFrom))
	if err != nil {
		return fmt.Errorf("%w: sender %q", ErrInvalidEmailAddress, s.cfg.EmailFrom)
	}

	encodedSubject := mime.QEncoding.Encode("utf-8", subject)

	if s.cfg.SMTPHost == "" {
		// Log-only simulation mode for tests and local development
		log.Printf("[EMAIL-SIMULATION] To: %s | From: %s | Subject: %s\n%s", parsedTo.Address, parsedFrom.Address, subject, plainBody)
		return nil
	}

	// 1. Build multipart/alternative body with plain text and HTML parts
	var bodyBuf bytes.Buffer
	mpWriter := multipart.NewWriter(&bodyBuf)

	// Plain text part (Base64)
	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	textHeader.Set("Content-Transfer-Encoding", "base64")
	partText, err := mpWriter.CreatePart(textHeader)
	if err != nil {
		return fmt.Errorf("failed to create text part: %w", err)
	}
	if _, err := partText.Write([]byte(base64.StdEncoding.EncodeToString([]byte(plainBody)))); err != nil {
		return fmt.Errorf("failed to write text part: %w", err)
	}

	// HTML part (Base64)
	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlHeader.Set("Content-Transfer-Encoding", "base64")
	partHTML, err := mpWriter.CreatePart(htmlHeader)
	if err != nil {
		return fmt.Errorf("failed to create html part: %w", err)
	}
	if _, err := partHTML.Write([]byte(base64.StdEncoding.EncodeToString([]byte(htmlBody)))); err != nil {
		return fmt.Errorf("failed to write html part: %w", err)
	}

	if err := mpWriter.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// 2. Build full RFC 5322 MIME message
	var msgBuf bytes.Buffer
	msgBuf.WriteString(fmt.Sprintf("From: %s\r\n", parsedFrom.Address))
	msgBuf.WriteString(fmt.Sprintf("To: %s\r\n", parsedTo.Address))
	msgBuf.WriteString(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	msgBuf.WriteString("MIME-Version: 1.0\r\n")
	msgBuf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", mpWriter.Boundary()))
	msgBuf.Write(bodyBuf.Bytes())

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" && s.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	}

	recipients := []string{parsedTo.Address}
	if err := smtp.SendMail(addr, auth, parsedFrom.Address, recipients, msgBuf.Bytes()); err != nil {
		log.Printf("[EMAIL-ERROR] Failed to send email to %s: %v", parsedTo.Address, err)
		return err
	}
	return nil
}
