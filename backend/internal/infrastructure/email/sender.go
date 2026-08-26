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

// sanitizeHeader cleans header values by strictly stripping CRLF and control characters (CWE-93).
func sanitizeHeader(input string) string {
	s := strings.ReplaceAll(input, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(s))
}

// sanitizeText removes CRLF and dangerous control characters from user-provided text.
func sanitizeText(input string) string {
	s := strings.ReplaceAll(input, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(s))
}

func (s *Service) SendVerificationEmail(ctx context.Context, toEmail, toName, token string) error {
	_ = ctx
	cleanName := sanitizeText(toName)
	if cleanName == "" {
		cleanName = "Usuário"
	}
	cleanToken := sanitizeHeader(token)
	verifyURL := fmt.Sprintf(
		"%s/verify-email?token=%s",
		strings.TrimRight(s.cfg.AppBaseURL, "/"),
		url.QueryEscape(cleanToken),
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
	cleanName := sanitizeText(toName)
	if cleanName == "" {
		cleanName = "Usuário"
	}
	cleanToken := sanitizeHeader(token)
	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		strings.TrimRight(s.cfg.AppBaseURL, "/"),
		url.QueryEscape(cleanToken),
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
	// Guard clauses against CRLF injection
	if strings.ContainsAny(toEmail, "\r\n") {
		return fmt.Errorf("%w: recipient contains newline characters", ErrInvalidEmailAddress)
	}
	if strings.ContainsAny(s.cfg.EmailFrom, "\r\n") {
		return fmt.Errorf("%w: sender contains newline characters", ErrInvalidEmailAddress)
	}
	if strings.ContainsAny(subject, "\r\n") {
		return errors.New("subject contains newline characters")
	}

	cleanTo := sanitizeHeader(toEmail)
	parsedTo, err := mail.ParseAddress(cleanTo)
	if err != nil {
		return fmt.Errorf("%w: recipient %q", ErrInvalidEmailAddress, toEmail)
	}
	safeToAddress := strings.ReplaceAll(strings.ReplaceAll(parsedTo.Address, "\r", ""), "\n", "")

	cleanFrom := sanitizeHeader(s.cfg.EmailFrom)
	parsedFrom, err := mail.ParseAddress(cleanFrom)
	if err != nil {
		return fmt.Errorf("%w: sender %q", ErrInvalidEmailAddress, s.cfg.EmailFrom)
	}
	safeFromAddress := strings.ReplaceAll(strings.ReplaceAll(parsedFrom.Address, "\r", ""), "\n", "")

	cleanSubject := sanitizeHeader(subject)
	safeSubject := strings.ReplaceAll(strings.ReplaceAll(cleanSubject, "\r", ""), "\n", "")
	encodedSubject := mime.QEncoding.Encode("utf-8", safeSubject)

	if s.cfg.SMTPHost == "" {
		// Log-only mode in local development / CI
		log.Printf("[EMAIL-SIMULATION] To: %s | From: %s | Subject: %s\n%s", safeToAddress, safeFromAddress, safeSubject, plainBody)
		return nil
	}

	// Build MIME multipart body using standard multipart.Writer
	var bodyBuf bytes.Buffer
	mpWriter := multipart.NewWriter(&bodyBuf)

	// 1. Text part (Base64 encoded)
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

	// 2. HTML part (Base64 encoded)
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

	// Build full RFC 5322 message
	var msgBuf bytes.Buffer
	msgBuf.WriteString("From: " + safeFromAddress + "\r\n")
	msgBuf.WriteString("To: " + safeToAddress + "\r\n")
	msgBuf.WriteString("Subject: " + encodedSubject + "\r\n")
	msgBuf.WriteString("MIME-Version: 1.0\r\n")
	msgBuf.WriteString("Content-Type: multipart/alternative; boundary=\"" + mpWriter.Boundary() + "\"\r\n\r\n")
	msgBuf.Write(bodyBuf.Bytes())

	safeHost := sanitizeHeader(s.cfg.SMTPHost)
	safePort := sanitizeHeader(s.cfg.SMTPPort)
	addr := fmt.Sprintf("%s:%s", safeHost, safePort)

	var auth smtp.Auth
	if s.cfg.SMTPUser != "" && s.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, safeHost)
	}

	recipients := []string{safeToAddress}
	if err := smtp.SendMail(addr, auth, safeFromAddress, recipients, msgBuf.Bytes()); err != nil {
		log.Printf("[EMAIL-ERROR] Failed to send email to %s: %v", safeToAddress, err)
		return err
	}
	return nil
}
