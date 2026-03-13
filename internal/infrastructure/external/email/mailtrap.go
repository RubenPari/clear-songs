package email

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/RubenPari/clear-songs/internal/domain/auth"
)

type mailtrapEmailService struct {
	host        string
	port        string
	username    string
	password    string
	from        string
	frontendURL string
}

func NewMailtrapEmailService() auth.EmailService {
	return &mailtrapEmailService{
		host:        os.Getenv("SMTP_HOST"),
		port:        os.Getenv("SMTP_PORT"),
		username:    os.Getenv("SMTP_USERNAME"),
		password:    os.Getenv("SMTP_PASSWORD"),
		from:        os.Getenv("SMTP_FROM"),
		frontendURL: os.Getenv("FRONTEND_URL"), // e.g., http://localhost:4200
	}
}

func (s *mailtrapEmailService) SendVerificationEmail(ctx context.Context, toEmail string, token string) error {
	subject := "Verify your email - Clear Songs"

	link := fmt.Sprintf("%s/confirm-email?token=%s", s.frontendURL, token)
	body := fmt.Sprintf("Welcome to Clear Songs!\n\nPlease click the link below to verify your email address:\n%s\n\nIf you did not request this, please ignore this email.", link)

	return s.sendEmail(toEmail, subject, body)
}

func (s *mailtrapEmailService) SendPasswordResetEmail(ctx context.Context, toEmail string, token string) error {
	subject := "Reset your password - Clear Songs"

	link := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)
	body := fmt.Sprintf("You requested a password reset for Clear Songs.\n\nPlease click the link below to reset your password:\n%s\n\nIf you did not request this, please ignore this email.", link)

	return s.sendEmail(toEmail, subject, body)
}

func (s *mailtrapEmailService) sendEmail(toEmail, subject, body string) error {
	if s.host == "" || s.port == "" || s.username == "" || s.password == "" {
		return fmt.Errorf("SMTP credentials are not fully configured")
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{toEmail}, msg.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
