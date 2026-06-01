package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/volo/volo-api/internal/model"
)

const (
	MaxBodySize      = 10 * 1024 // 10KB max request body
	MaxTranscriptLen = 500       // Max voice transcript length
	MaxMessageLen    = 2000      // Max chat message length
	MaxEmailLen      = 255
	MaxPasswordLen   = 128
	MinPasswordLen   = 8
	MaxDeviceNameLen = 100
)

// ValidateBody limits request body size to prevent abuse.
func ValidateBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil && r.ContentLength > int64(MaxBodySize) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(model.Response{
				Error: &model.ErrorResponse{
					Code:    "BODY_TOO_LARGE",
					Message: "Request body exceeds maximum size (10KB).",
				},
			})
			return
		}

		// Wrap body with a size limiter
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, int64(MaxBodySize))
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateTranscript checks transcript field constraints.
func ValidateTranscript(transcript string) *model.ErrorResponse {
	if strings.TrimSpace(transcript) == "" {
		return &model.ErrorResponse{
			Code:    "EMPTY_TRANSCRIPT",
			Message: "Transcript cannot be empty.",
		}
	}
	if len(transcript) > MaxTranscriptLen {
		return &model.ErrorResponse{
			Code:    "TRANSCRIPT_TOO_LONG",
			Message: "Transcript exceeds maximum length (500 characters).",
		}
	}
	return nil
}

// ValidateChatMessage checks chat message constraints.
func ValidateChatMessage(message string) *model.ErrorResponse {
	if strings.TrimSpace(message) == "" {
		return &model.ErrorResponse{
			Code:    "EMPTY_MESSAGE",
			Message: "Message cannot be empty.",
		}
	}
	if len(message) > MaxMessageLen {
		return &model.ErrorResponse{
			Code:    "MESSAGE_TOO_LONG",
			Message: "Message exceeds maximum length (2000 characters).",
		}
	}
	return nil
}

// ValidateEmail checks email format (basic).
func ValidateEmail(email string) *model.ErrorResponse {
	email = strings.TrimSpace(email)
	if email == "" {
		return &model.ErrorResponse{
			Code:    "MISSING_EMAIL",
			Message: "Email is required.",
		}
	}
	if len(email) > MaxEmailLen {
		return &model.ErrorResponse{
			Code:    "EMAIL_TOO_LONG",
			Message: "Email exceeds maximum length.",
		}
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return &model.ErrorResponse{
			Code:    "INVALID_EMAIL",
			Message: "Please provide a valid email address.",
		}
	}
	return nil
}

// ValidatePassword checks password constraints.
func ValidatePassword(password string) *model.ErrorResponse {
	if len(password) < MinPasswordLen {
		return &model.ErrorResponse{
			Code:    "WEAK_PASSWORD",
			Message: "Password must be at least 8 characters.",
		}
	}
	if len(password) > MaxPasswordLen {
		return &model.ErrorResponse{
			Code:    "PASSWORD_TOO_LONG",
			Message: "Password exceeds maximum length.",
		}
	}
	return nil
}

// DrainBody reads and discards the remaining body to allow connection reuse.
func DrainBody(body io.ReadCloser) {
	if body != nil {
		io.Copy(io.Discard, body)
		body.Close()
	}
}
