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
	"net/smtp"
	"net/textproto"
	"net/url"
	"regexp"
	"strings"

	emailport "cvmc/internal/application/ports/email"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	ErrInvalidEmailAddress = errors.New("invalid email address")

	verificationTmpl  = htmltemplate.Must(htmltemplate.ParseFS(templateFS, "templates/verification.html"))
	passwordResetTmpl = htmltemplate.Must(htmltemplate.ParseFS(templateFS, "templates/password_reset.html"))

	// emailRegex validates that an address contains only safe ASCII characters.
	// This acts as a CodeQL-recognized sanitizer barrier: after matching,
	// the string is proven to contain no CRLF, control chars, or injection vectors.
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
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

// validateEmail checks the address against a strict ASCII-only regex.
// CodeQL recognizes regexp.MatchString / Regexp.MatchString as a taint sanitizer.
func validateEmail(address string) (string, error) {
	clean := strings.TrimSpace(address)
	if !emailRegex.MatchString(clean) {
		return "", fmt.Errorf("%w: %q", ErrInvalidEmailAddress, address)
	}
	return clean, nil
}

func sanitizeHeaderValue(value string) (string, error) {
	clean := strings.TrimSpace(value)
	if strings.ContainsAny(clean, "\r\n") {
		return "", errors.New("invalid header value")
	}
	return clean, nil
}

func (s *Service) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	_ = ctx

	// token is server-generated (crypto/rand), not user input — safe by construction.
	safeToken := url.QueryEscape(token)
	verifyURL := strings.TrimRight(s.cfg.AppBaseURL, "/") + "/verify-email?token=" + safeToken

	subject := "CVMC - Confirmação de E-mail"

	var htmlBuf bytes.Buffer
	if err := verificationTmpl.Execute(&htmlBuf, struct{ VerifyURL string }{VerifyURL: verifyURL}); err != nil {
		return fmt.Errorf("failed to render verification email template: %w", err)
	}

	plainBody := "Olá!\n\nObrigado por se cadastrar no CVMC.\n\nPara validar seu e-mail, acesse:\n" + verifyURL + "\n\nEste link é válido por 24 horas.\n\nSe você não criou uma conta, ignore este e-mail."

	return s.sendEmail(toEmail, subject, plainBody, htmlBuf.String())
}

func (s *Service) SendPasswordResetEmail(ctx context.Context, toEmail, token string) error {
	_ = ctx

	safeToken := url.QueryEscape(token)
	resetURL := strings.TrimRight(s.cfg.AppBaseURL, "/") + "/reset-password?token=" + safeToken

	subject := "CVMC - Recuperação de Senha"

	var htmlBuf bytes.Buffer
	if err := passwordResetTmpl.Execute(&htmlBuf, struct{ ResetURL string }{ResetURL: resetURL}); err != nil {
		return fmt.Errorf("failed to render password reset email template: %w", err)
	}

	plainBody := "Olá!\n\nPara redefinir sua senha no CVMC, acesse:\n" + resetURL + "\n\nEste link é válido por 30 minutos.\n\nSe você não solicitou, ignore este e-mail."

	return s.sendEmail(toEmail, subject, plainBody, htmlBuf.String())
}

// sendEmail validates addresses with regex (CodeQL sanitizer), builds a safe MIME message,
// and dispatches via smtp.SendMail.
func (s *Service) sendEmail(toEmail, subject, plainBody, htmlBody string) error {
	// Validate recipient with strict regex — this is the CodeQL taint barrier.
	validTo, err := validateEmail(toEmail)
	if err != nil {
		return err
	}

	// From address comes from server config, not user input.
	validFrom, err := validateEmail(s.cfg.EmailFrom)
	if err != nil {
		return err
	}

	safeTo, err := sanitizeHeaderValue(validTo)
	if err != nil {
		return err
	}
	safeFrom, err := sanitizeHeaderValue(validFrom)
	if err != nil {
		return err
	}
	safeSubject, err := sanitizeHeaderValue(subject)
	if err != nil {
		return err
	}

	encodedSubject := mime.QEncoding.Encode("utf-8", safeSubject)

	if s.cfg.SMTPHost == "" {
		log.Printf("[EMAIL-SIMULATION] To: %s | Subject: %s\n%s", safeTo, safeSubject, plainBody)
		return nil
	}

	// Build MIME multipart/alternative body
	var bodyBuf bytes.Buffer
	mpWriter := multipart.NewWriter(&bodyBuf)

	textHdr := make(textproto.MIMEHeader)
	textHdr.Set("Content-Type", "text/plain; charset=UTF-8")
	textHdr.Set("Content-Transfer-Encoding", "base64")
	pw, err := mpWriter.CreatePart(textHdr)
	if err != nil {
		return fmt.Errorf("failed to create text part: %w", err)
	}
	_, _ = pw.Write([]byte(base64.StdEncoding.EncodeToString([]byte(plainBody))))

	htmlHdr := make(textproto.MIMEHeader)
	htmlHdr.Set("Content-Type", "text/html; charset=UTF-8")
	htmlHdr.Set("Content-Transfer-Encoding", "base64")
	hw, err := mpWriter.CreatePart(htmlHdr)
	if err != nil {
		return fmt.Errorf("failed to create html part: %w", err)
	}
	_, _ = hw.Write([]byte(base64.StdEncoding.EncodeToString([]byte(htmlBody))))

	_ = mpWriter.Close()

	// Assemble RFC 5322 message — only sanitized values enter the headers.
	var msg bytes.Buffer
	msg.WriteString("From: " + safeFrom + "\r\n")
	msg.WriteString("To: " + safeTo + "\r\n")
	msg.WriteString("Subject: " + encodedSubject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: multipart/alternative; boundary=\"" + mpWriter.Boundary() + "\"\r\n\r\n")
	msg.Write(bodyBuf.Bytes())

	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" && s.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	}

	if err := smtp.SendMail(addr, auth, safeFrom, []string{safeTo}, msg.Bytes()); err != nil {
		log.Printf("[EMAIL-ERROR] Failed to send email to %s: %v", safeTo, err)
		return err
	}
	return nil
}
