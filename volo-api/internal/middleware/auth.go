package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/volo/volo-api/internal/model"
)

const UserIDKey contextKey = "user_id"

// AuthService is the interface the auth middleware needs.
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (userID string, err error)
}

// Auth validates JWT tokens and injects user_id into context.
func Auth(authSvc AuthService, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				unauthorized(w, "Missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				unauthorized(w, "Invalid authorization format. Use: Bearer <token>")
				return
			}

			token := parts[1]
			userID, err := authSvc.ValidateToken(r.Context(), token)
			if err != nil {
				unauthorized(w, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the authenticated user ID from context.
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(model.Response{
		Error: &model.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: msg,
		},
	})
}
