package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const resendURL = "https://api.resend.com/emails"

// EmailService sends transactional emails via Resend.
type EmailService struct {
	apiKey   string
	fromAddr string
	client   *http.Client
}

func NewEmailService(apiKey string) *EmailService {
	fromAddr := "Volo <noreply@volo.app>"
	return &EmailService{
		apiKey:   apiKey,
		fromAddr: fromAddr,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IsConfigured returns true if the email service has an API key.
func (s *EmailService) IsConfigured() bool {
	return s.apiKey != ""
}

// SendPasswordReset sends a password reset email.
func (s *EmailService) SendPasswordReset(to, resetToken, resetURL string) error {
	link := fmt.Sprintf("%s?token=%s", resetURL, resetToken)

	subject := "Reset your Volo password"
	html := fmt.Sprintf(`
		<div style="font-family: -apple-system, system-ui, sans-serif; max-width: 480px; margin: 0 auto; padding: 40px 20px;">
			<h2 style="color: #1a1a2e; margin-bottom: 16px;">Reset your password</h2>
			<p style="color: #4a4a4a; line-height: 1.6;">
				You requested a password reset for your Volo account. Click the button below to set a new password.
			</p>
			<a href="%s" style="display: inline-block; background: #3b82f6; color: white; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 500; margin: 24px 0;">
				Reset Password
			</a>
			<p style="color: #888; font-size: 13px; line-height: 1.5;">
				This link expires in 1 hour. If you didn't request this, you can safely ignore this email.
			</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;" />
			<p style="color: #aaa; font-size: 12px;">Volo — Voice-first browser assistant</p>
		</div>
	`, link)

	return s.send(to, subject, html)
}

func (s *EmailService) send(to, subject, html string) error {
	if !s.IsConfigured() {
		slog.Warn("email service not configured, skipping send", "to", to, "subject", subject)
		return nil
	}

	payload := map[string]any{
		"from":    s.fromAddr,
		"to":      []string{to},
		"subject": subject,
		"html":    html,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", resendURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("email.send → create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("email.send → request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("email.send → resend returned %d: %s", resp.StatusCode, string(respBody))
	}

	slog.Info("email sent", "to", to, "subject", subject)
	return nil
}
