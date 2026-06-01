package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/volo/volo-api/internal/model"
)

// Recoverer catches panics and returns a 500 JSON response.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"request_id", GetRequestID(r.Context()),
					"error", err,
					"stack", string(debug.Stack()),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(model.Response{
					Error: &model.ErrorResponse{
						Code:          "INTERNAL_ERROR",
						Message:       "An unexpected error occurred. Please try again.",
						InternalError: "panic recovered — check server logs",
					},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
