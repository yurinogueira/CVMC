package email

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"
	"unicode"

	emailport "cvmc/internal/application/ports/email"
)

var (
	ErrInvalidEmailAddress = errors.New("invalid email address")
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
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' || unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(input))
}

// sanitizeContent removes dangerous control characters from body variables.
func sanitizeContent(input string) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return ' '
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, strings.TrimSpace(input))
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
	body := fmt.Sprintf(
		"Olá, %s!\n\nObrigado por se cadastrar no CVMC (Como Vai Meu Carro).\n\nPara validar seu e-mail e liberar o cadastro de veículos, acesse o link abaixo:\n%s\n\nEste link é válido por 24 horas.\n\nSe você não criou uma conta no CVMC, ignore este e-mail.",
		cleanName,
		verifyURL,
	)

	return s.send(toEmail, subject, body)
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
	body := fmt.Sprintf(
		"Olá, %s!\n\nRecebemos uma solicitação para redefinir a senha da sua conta no CVMC.\n\nPara criar uma nova senha, acesse o link abaixo:\n%s\n\nEste link é válido por 30 minutos.\n\nSe você não solicitou a redefinição de senha, ignore este e-mail com segurança.",
		cleanName,
		resetURL,
	)

	return s.send(toEmail, subject, body)
}

func (s *Service) send(toEmail, subject, body string) error {
	cleanTo := sanitizeHeader(toEmail)
	parsedTo, err := mail.ParseAddress(cleanTo)
	if err != nil {
		return fmt.Errorf("%w: recipient %q", ErrInvalidEmailAddress, toEmail)
	}

	cleanFrom := sanitizeHeader(s.cfg.EmailFrom)
	parsedFrom, err := mail.ParseAddress(cleanFrom)
	if err != nil {
		return fmt.Errorf("%w: sender %q", ErrInvalidEmailAddress, s.cfg.EmailFrom)
	}

	cleanSubject := sanitizeHeader(subject)
	encodedSubject := mime.QEncoding.Encode("utf-8", cleanSubject)

	if s.cfg.SMTPHost == "" {
		// Log-only mode in local development / CI
		log.Printf("[EMAIL-SIMULATION] To: %s | From: %s | Subject: %s\n%s", parsedTo.Address, parsedFrom.Address, cleanSubject, body)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", sanitizeHeader(s.cfg.SMTPHost), sanitizeHeader(s.cfg.SMTPPort))
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" && s.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, sanitizeHeader(s.cfg.SMTPHost))
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		parsedFrom.Address,
		parsedTo.Address,
		encodedSubject,
		body,
	))

	if err := smtp.SendMail(addr, auth, parsedFrom.Address, []string{parsedTo.Address}, msg); err != nil {
		log.Printf("[EMAIL-ERROR] Failed to send email to %s: %v", parsedTo.Address, err)
		return err
	}
	return nil
}
