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

func (h *Handlers) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req model.ForgotPasswordRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if errResp := middleware.ValidateEmail(req.Email); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}

	rawToken, err := h.services.Auth.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "FORGOT_PASSWORD_FAILED", "Something went wrong. Please try again.", err.Error())
		return
	}

	// Send email (if token was generated — won't be if email doesn't exist or rate limited)
	if rawToken != "" {
		go func() {
			if err := h.services.Email.SendPasswordReset(req.Email, rawToken, h.cfg.ResetPasswordURL); err != nil {
				middleware.CaptureError(err, map[string]string{"action": "send_password_reset"})
			}
		}()
	}

	// Always return 200 — don't reveal if email exists
	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"message": "If an account with that email exists, a reset link has been sent."},
	})
}

func (h *Handlers) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req model.ResetPasswordRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Token == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_TOKEN", "Reset token is required.", "")
		return
	}

	if errResp := middleware.ValidatePassword(req.Password); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}

	if err := h.services.Auth.ResetPassword(r.Context(), req.Token, req.Password); err != nil {
		h.respondError(w, http.StatusBadRequest, "RESET_FAILED", "Invalid or expired reset link. Please request a new one.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"message": "Password reset successful. Please log in with your new password."},
	})
}

func (h *Handlers) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req model.VerifyEmailRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Token == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_TOKEN", "Verification token is required.", "")
		return
	}

	if err := h.services.Auth.VerifyEmail(r.Context(), req.Token); err != nil {
		h.respondError(w, http.StatusBadRequest, "VERIFICATION_FAILED", "Invalid or expired verification link.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"message": "Email verified successfully."},
	})
}

func (h *Handlers) ResendVerification(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := h.services.Auth.GetUserByID(r.Context(), userID)
	if err != nil || user.Email == nil {
		h.respondError(w, http.StatusBadRequest, "NO_EMAIL", "No email associated with this account.", "")
		return
	}

	if user.EmailVerifiedAt != nil {
		h.respondError(w, http.StatusBadRequest, "ALREADY_VERIFIED", "Email is already verified.", "")
		return
	}

	rawToken, err := h.services.Auth.SendVerificationEmail(r.Context(), userID, *user.Email)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "VERIFICATION_SEND_FAILED", "Failed to send verification email.", err.Error())
		return
	}

	if rawToken != "" {
		go func() {
			if err := h.services.Email.SendEmailVerification(*user.Email, rawToken, h.cfg.VerifyEmailURL); err != nil {
				middleware.CaptureError(err, map[string]string{"action": "send_verification_email"})
			}
		}()
	}

	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"message": "Verification email sent."},
	})
}

func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req model.ChangePasswordRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.CurrentPassword == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_CURRENT_PASSWORD", "Current password is required.", "")
		return
	}

	if errResp := middleware.ValidatePassword(req.NewPassword); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}

	if req.CurrentPassword == req.NewPassword {
		h.respondError(w, http.StatusBadRequest, "SAME_PASSWORD", "New password must be different from current password.", "")
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.services.Auth.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		h.respondError(w, http.StatusBadRequest, "CHANGE_PASSWORD_FAILED", "Failed to change password. Check your current password.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"message": "Password changed successfully."},
	})
}

func (h *Handlers) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var req model.DeleteAccountRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	if req.Password == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_PASSWORD", "Password is required to confirm account deletion.", "")
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.services.Auth.DeleteAccount(r.Context(), userID, req.Password); err != nil {
		h.respondError(w, http.StatusBadRequest, "DELETE_FAILED", "Failed to delete account. Check your password.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]string{"message": "Account deleted. We're sorry to see you go."},
	})
}
