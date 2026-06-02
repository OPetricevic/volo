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
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.DeviceID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DEVICE_ID", "Device ID is required.", "")
		return
	}

	// Validate device name length
	if len(req.DeviceName) > middleware.MaxDeviceNameLen {
		h.respondError(w, http.StatusBadRequest, "DEVICE_NAME_TOO_LONG", "Device name exceeds maximum length.", "")
		return
	}

	result, err := h.services.Auth.RegisterDevice(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register device. Please try again.", err.Error())
		return
	}

	middleware.RecordRegistration()
	h.respond(w, http.StatusCreated, model.Response{Data: result})
}

func (h *Handlers) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req model.GoogleAuthRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Code == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_CODE", "Authorization code is required.", "")
		return
	}

	if req.DeviceID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DEVICE_ID", "Device ID is required.", "")
		return
	}

	// Check if Google OAuth is configured
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleClientSecret == "" {
		h.respondError(w, http.StatusServiceUnavailable, "GOOGLE_NOT_CONFIGURED", "Google sign-in is not configured on this server.", "")
		return
	}

	result, err := h.services.Auth.GoogleAuth(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "GOOGLE_AUTH_FAILED", "Google sign-in failed. Please try again.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: result})
}

func (h *Handlers) RegisterEmail(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterEmailRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	// Use proper validation functions
	if errResp := middleware.ValidateEmail(req.Email); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}
	if errResp := middleware.ValidatePassword(req.Password); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}

	userID := middleware.GetUserID(r.Context()) // may be empty if not authenticated
	result, err := h.services.Auth.RegisterEmail(r.Context(), req, userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to create account. Email may already be in use.", err.Error())
		return
	}

	h.respond(w, http.StatusCreated, model.Response{Data: result})
}

func (h *Handlers) LoginEmail(w http.ResponseWriter, r *http.Request) {
	var req model.LoginEmailRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required.", "")
		return
	}

	result, err := h.services.Auth.LoginEmail(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password.", "")
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: result})
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_TOKEN", "No token provided.", "")
		return
	}

	if err := h.services.Auth.Logout(r.Context(), token); err != nil {
		h.respondError(w, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "logged_out"}})
}

func (h *Handlers) LogoutAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if err := h.services.Auth.LogoutAll(r.Context(), userID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout all devices.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "all_sessions_revoked"}})
}

func (h *Handlers) UnlinkDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_DEVICE_ID", "Device ID is required.", "")
		return
	}

	userID := middleware.GetUserID(r.Context())

	// Verify the device belongs to the authenticated user
	device, err := h.services.Auth.GetUserDevice(r.Context(), userID, deviceID)
	if err != nil || device == nil {
		h.respondError(w, http.StatusNotFound, "DEVICE_NOT_FOUND", "Device not found or does not belong to you.", "")
		return
	}

	if err := h.services.Auth.UnlinkDevice(r.Context(), userID, deviceID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "UNLINK_FAILED", "Failed to unlink device.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "device_unlinked"}})
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
