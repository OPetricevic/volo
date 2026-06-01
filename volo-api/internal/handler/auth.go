package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/volo/volo-api/internal/middleware"
	"github.com/volo/volo-api/internal/model"
)

func (h *Handlers) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterDeviceRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.DeviceID == "" {
		respondError(w, http.StatusBadRequest, "MISSING_DEVICE_ID", "Device ID is required.", "")
		return
	}

	result, err := h.services.Auth.RegisterDevice(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register device. Please try again.", err.Error())
		return
	}

	middleware.RecordRegistration()
	respond(w, http.StatusCreated, model.Response{Data: result})
}

func (h *Handlers) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req model.GoogleAuthRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Code == "" {
		respondError(w, http.StatusBadRequest, "MISSING_CODE", "Authorization code is required.", "")
		return
	}

	if req.DeviceID == "" {
		respondError(w, http.StatusBadRequest, "MISSING_DEVICE_ID", "Device ID is required.", "")
		return
	}

	// Check if Google OAuth is configured
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleClientSecret == "" {
		respondError(w, http.StatusServiceUnavailable, "GOOGLE_NOT_CONFIGURED", "Google sign-in is not configured on this server.", "")
		return
	}

	result, err := h.services.Auth.GoogleAuth(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "GOOGLE_AUTH_FAILED", "Google sign-in failed. Please try again.", err.Error())
		return
	}

	respond(w, http.StatusOK, model.Response{Data: result})
}

func (h *Handlers) RegisterEmail(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterEmailRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required.", "")
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 8 characters.", "")
		return
	}

	userID := middleware.GetUserID(r.Context()) // may be empty if not authenticated
	result, err := h.services.Auth.RegisterEmail(r.Context(), req, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to create account. Email may already be in use.", err.Error())
		return
	}

	respond(w, http.StatusCreated, model.Response{Data: result})
}

func (h *Handlers) LoginEmail(w http.ResponseWriter, r *http.Request) {
	var req model.LoginEmailRequest
	if err := decode(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required.", "")
		return
	}

	result, err := h.services.Auth.LoginEmail(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password.", err.Error())
		return
	}

	respond(w, http.StatusOK, model.Response{Data: result})
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		respondError(w, http.StatusBadRequest, "MISSING_TOKEN", "No token provided.", "")
		return
	}

	if err := h.services.Auth.Logout(r.Context(), token); err != nil {
		respondError(w, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout.", err.Error())
		return
	}

	respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "logged_out"}})
}

func (h *Handlers) LogoutAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if err := h.services.Auth.LogoutAll(r.Context(), userID); err != nil {
		respondError(w, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout all devices.", err.Error())
		return
	}

	respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "all_sessions_revoked"}})
}

func (h *Handlers) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		respondError(w, http.StatusBadRequest, "MISSING_DEVICE_ID", "Device ID is required.", "")
		return
	}

	// TODO: Verify the device belongs to the authenticated user
	respondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Device unlinking is not yet available.", "")
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
