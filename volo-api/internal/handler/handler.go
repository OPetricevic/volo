package handler

import (
	"encoding/json"
	"net/http"

	"github.com/volo/volo-api/internal/config"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/service"
)

// Handlers holds all HTTP handlers.
type Handlers struct {
	services *service.Services
	cfg      *config.Config
}

// New creates the handlers.
func New(services *service.Services, cfg *config.Config) *Handlers {
	return &Handlers{services: services, cfg: cfg}
}

// Health returns a simple health check.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"status": "ok"},
	})
}

// --- Helpers ---

func (h *Handlers) respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func (h *Handlers) respondError(w http.ResponseWriter, status int, code, message, internalErr string) {
	errResp := &model.ErrorResponse{
		Code:    code,
		Message: message,
	}
	// Only include internal error details in non-production environments
	if h.cfg != nil && h.cfg.Environment == "production" {
		// Suppress internal errors in production
	} else {
		errResp.InternalError = internalErr
	}
	h.respond(w, status, model.Response{Error: errResp})
}

// respond is a standalone helper for contexts without Handlers (e.g. Health)
func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, status int, code, message, internalErr string) {
	respond(w, status, model.Response{
		Error: &model.ErrorResponse{
			Code:          code,
			Message:       message,
			InternalError: internalErr,
		},
	})
}

func decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
