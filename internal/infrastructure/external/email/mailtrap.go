package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/RubenPari/clear-songs/internal/domain/auth"
)

type mailtrapEmailService struct {
	apiToken    string
	from        string
	frontendURL string
}

func NewMailtrapEmailService() auth.EmailService {
	return &mailtrapEmailService{
		apiToken:    os.Getenv("MAILTRAP_API_TOKEN"),
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

// Mailtrap structures based on the API documentation
type Recipient struct {
	Email string `json:"email"`
}

type Sender struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailtrapRequest struct {
	From    Sender      `json:"from"`
	To      []Recipient `json:"to"`
	Subject string      `json:"subject"`
	Text    string      `json:"text,omitempty"`
}

func (s *mailtrapEmailService) sendEmail(toEmail, subject, body string) error {
	if s.apiToken == "" {
		return fmt.Errorf("MAILTRAP_API_TOKEN is not configured")
	}

	url := "https://send.api.mailtrap.io/api/send"

	emailBody := MailtrapRequest{
		From: Sender{
			Email: s.from,
			Name:  "Clear Songs",
		},
		To: []Recipient{
			{Email: toEmail},
		},
		Subject: subject,
		Text:    body,
	}

	jsonData, err := json.Marshal(emailBody)
	if err != nil {
		return fmt.Errorf("error marshaling email data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating email request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var respBody bytes.Buffer
		respBody.ReadFrom(resp.Body)
		return fmt.Errorf("failed to send email, status code: %d, response: %s", resp.StatusCode, respBody.String())
	}

	return nil
}
