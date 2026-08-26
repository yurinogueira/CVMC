package email

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"log"
	"mime"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"
	"unicode"

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

// sanitizeHeader removes CRLF and control characters to strictly prevent SMTP header injection (CWE-93).
func sanitizeHeader(input string) string {
	noCR := strings.ReplaceAll(input, "\r", "")
	noLF := strings.ReplaceAll(noCR, "\n", "")
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(noLF))
}

// sanitizeContent removes dangerous control characters from body variables.
func sanitizeContent(input string) string {
	noCR := strings.ReplaceAll(input, "\r", " ")
	noLF := strings.ReplaceAll(noCR, "\n", " ")
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(noLF))
}

// encodeBase64Body formats a payload with standard RFC 2045 line wrapping (76 chars) for safe MIME transport.
func encodeBase64Body(content string) string {
	raw := base64.StdEncoding.EncodeToString([]byte(content))
	var buf strings.Builder
	for len(raw) > 76 {
		buf.WriteString(raw[:76])
		buf.WriteString("\r\n")
		raw = raw[76:]
	}
	buf.WriteString(raw)
	return buf.String()
}

func generateBoundary() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "cvmc_multipart_boundary_default"
	}
	return "cvmc_bnd_" + hex.EncodeToString(b)
}

func (s *Service) SendVerificationEmail(ctx context.Context, toEmail, toName, token string) error {
	_ = ctx
	cleanName := sanitizeContent(toName)
	if cleanName == "" {
		cleanName = "Usuário"
	}
	verifyURL := fmt.Sprintf(
		"%s/verify-email?token=%s",
		strings.TrimRight(s.cfg.AppBaseURL, "/"),
		url.QueryEscape(sanitizeHeader(token)),
	)
	subject := "CVMC - Confirmação de E-mail"

	var htmlBuf bytes.Buffer
	data := struct {
		Name      string
		VerifyURL string
	}{
		Name:      cleanName,
		VerifyURL: verifyURL,
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
	cleanName := sanitizeContent(toName)
	if cleanName == "" {
		cleanName = "Usuário"
	}
	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		strings.TrimRight(s.cfg.AppBaseURL, "/"),
		url.QueryEscape(sanitizeHeader(token)),
	)
	subject := "CVMC - Recuperação de Senha"

	var htmlBuf bytes.Buffer
	data := struct {
		Name     string
		ResetURL string
	}{
		Name:     cleanName,
		ResetURL: resetURL,
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
	if strings.ContainsAny(toEmail, "\r\n") {
		return fmt.Errorf("%w: recipient contains newline characters", ErrInvalidEmailAddress)
	}
	cleanTo := sanitizeHeader(toEmail)
	parsedTo, err := mail.ParseAddress(cleanTo)
	if err != nil {
		return fmt.Errorf("%w: recipient %q", ErrInvalidEmailAddress, toEmail)
	}
	safeToAddress := strings.ReplaceAll(strings.ReplaceAll(parsedTo.Address, "\r", ""), "\n", "")

	if strings.ContainsAny(s.cfg.EmailFrom, "\r\n") {
		return fmt.Errorf("%w: sender contains newline characters", ErrInvalidEmailAddress)
	}
	cleanFrom := sanitizeHeader(s.cfg.EmailFrom)
	parsedFrom, err := mail.ParseAddress(cleanFrom)
	if err != nil {
		return fmt.Errorf("%w: sender %q", ErrInvalidEmailAddress, s.cfg.EmailFrom)
	}
	safeFromAddress := strings.ReplaceAll(strings.ReplaceAll(parsedFrom.Address, "\r", ""), "\n", "")

	if strings.ContainsAny(subject, "\r\n") {
		return errors.New("subject contains newline characters")
	}
	cleanSubject := sanitizeHeader(subject)
	encodedSubject := mime.QEncoding.Encode("utf-8", cleanSubject)

	safePlainBody := sanitizeContent(plainBody)
	safeHTMLBody := sanitizeContent(htmlBody)

	if s.cfg.SMTPHost == "" {
		// Log-only mode in local development / CI
		log.Printf("[EMAIL-SIMULATION] To: %s | From: %s | Subject: %s\n%s", safeToAddress, safeFromAddress, cleanSubject, safePlainBody)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", sanitizeHeader(s.cfg.SMTPHost), sanitizeHeader(s.cfg.SMTPPort))
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" && s.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, sanitizeHeader(s.cfg.SMTPHost))
	}

	boundary := generateBoundary()
	b64Plain := encodeBase64Body(safePlainBody)
	b64HTML := encodeBase64Body(safeHTMLBody)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"Content-Transfer-Encoding: base64\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"Content-Transfer-Encoding: base64\r\n\r\n"+
			"%s\r\n\r\n"+
			"--%s--\r\n",
		safeFromAddress,
		safeToAddress,
		encodedSubject,
		boundary,
		boundary,
		b64Plain,
		boundary,
		b64HTML,
		boundary,
	))

	recipients := []string{safeToAddress}
	if err := smtp.SendMail(addr, auth, safeFromAddress, recipients, msg); err != nil {
		log.Printf("[EMAIL-ERROR] Failed to send email to %s: %v", safeToAddress, err)
		return err
	}
	return nil
}
