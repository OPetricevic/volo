package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/volo/volo-api/internal/model"
)

// Mock auth service
type mockAuthService struct {
	validateFunc func(ctx context.Context, token string) (string, error)
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	return m.validateFunc(ctx, token)
}

func TestAuth_ValidToken(t *testing.T) {
	mock := &mockAuthService{
		validateFunc: func(ctx context.Context, token string) (string, error) {
			if token == "valid-token" {
				return "user-123", nil
			}
			return "", errors.New("invalid")
		},
	}

	handler := Auth(mock, "secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID != "user-123" {
			t.Errorf("userID = %q, want user-123", userID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	mock := &mockAuthService{
		validateFunc: func(ctx context.Context, token string) (string, error) {
			return "", errors.New("invalid")
		},
	}

	handler := Auth(mock, "secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "UNAUTHORIZED" {
		t.Error("expected UNAUTHORIZED error")
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	mock := &mockAuthService{
		validateFunc: func(ctx context.Context, token string) (string, error) {
			return "", errors.New("invalid")
		},
	}

	handler := Auth(mock, "secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	mock := &mockAuthService{
		validateFunc: func(ctx context.Context, token string) (string, error) {
			return "", errors.New("expired")
		},
	}

	handler := Auth(mock, "secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
}

func TestGetUserID_MissingContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	id := GetUserID(req.Context())
	if id != "" {
		t.Errorf("expected empty string, got %q", id)
	}
}
