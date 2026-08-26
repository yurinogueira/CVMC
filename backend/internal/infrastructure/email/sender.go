package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	emailport "cvmc/internal/application/ports/email"
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
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", strings.TrimRight(s.cfg.AppBaseURL, "/"), token)
	subject := "CVMC - Confirmação de E-mail"
	body := fmt.Sprintf("Olá, %s!\n\nObrigado por se cadastrar no CVMC (Como Vai Meu Carro).\n\nPara validar seu e-mail e liberar o cadastro de veículos, acesse o link abaixo:\n%s\n\nEste link é válido por 24 horas.\n\nSe você não criou uma conta no CVMC, ignore este e-mail.", toName, verifyURL)

	return s.send(toEmail, subject, body)
}

func (s *Service) SendPasswordResetEmail(ctx context.Context, toEmail, toName, token string) error {
	_ = ctx
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimRight(s.cfg.AppBaseURL, "/"), token)
	subject := "CVMC - Recuperação de Senha"
	body := fmt.Sprintf("Olá, %s!\n\nRecebemos uma solicitação para redefinir a senha da sua conta no CVMC.\n\nPara criar uma nova senha, acesse o link abaixo:\n%s\n\nEste link é válido por 30 minutos.\n\nSe você não solicitou a redefinição de senha, ignore este e-mail com segurança.", toName, resetURL)

	return s.send(toEmail, subject, body)
}

func (s *Service) send(toEmail, subject, body string) error {
	if s.cfg.SMTPHost == "" {
		// Log-only mode in local development / CI
		log.Printf("[EMAIL-SIMULATION] To: %s | Subject: %s\n%s", toEmail, subject, body)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	var auth smtp.Auth
	if s.cfg.SMTPUser != "" && s.cfg.SMTPPass != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	}

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.EmailFrom, toEmail, subject, body))

	err := smtp.SendMail(addr, auth, s.cfg.EmailFrom, []string{toEmail}, msg)
	if err != nil {
		log.Printf("[EMAIL-ERROR] Failed to send email to %s: %v", toEmail, err)
		return err
	}
	return nil
}
