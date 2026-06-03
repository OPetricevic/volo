package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/volo/volo-api/internal/config"
	"github.com/volo/volo-api/internal/model"
)

func TestRegisterDevice_MissingDeviceID(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"device_id":"","device_name":"Chrome","platform":"extension"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/device", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.RegisterDevice(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_DEVICE_ID" {
		t.Errorf("expected MISSING_DEVICE_ID, got %v", resp.Error)
	}
}

func TestRegisterDevice_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{broken`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/device", body)
	rec := httptest.NewRecorder()

	h.RegisterDevice(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestRegisterEmail_MissingFields(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"email":"","password":"","device_id":"dev-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.RegisterEmail(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_EMAIL" {
		t.Errorf("expected MISSING_EMAIL, got %v", resp.Error)
	}
}

func TestRegisterEmail_WeakPassword(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"email":"test@test.com","password":"short","device_id":"dev-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.RegisterEmail(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "WEAK_PASSWORD" {
		t.Errorf("expected WEAK_PASSWORD, got %v", resp.Error)
	}
}

func TestLoginEmail_MissingFields(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"email":"","password":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.LoginEmail(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestLogout_MissingToken(t *testing.T) {
	h := &Handlers{}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	// No Authorization header
	rec := httptest.NewRecorder()

	h.Logout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestGoogleAuth_MissingCode(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"code":"","device_id":"dev-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.GoogleAuth(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_CODE" {
		t.Errorf("expected MISSING_CODE, got %v", resp.Error)
	}
}

func TestGoogleAuth_MissingDeviceID(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"code":"auth-code-123","device_id":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.GoogleAuth(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_DEVICE_ID" {
		t.Errorf("expected MISSING_DEVICE_ID, got %v", resp.Error)
	}
}

func TestGoogleAuth_NotConfigured(t *testing.T) {
	h := &Handlers{cfg: &config.Config{GoogleClientID: "", GoogleClientSecret: ""}}

	body := bytes.NewBufferString(`{"code":"auth-code-123","device_id":"dev-1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.GoogleAuth(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "GOOGLE_NOT_CONFIGURED" {
		t.Errorf("expected GOOGLE_NOT_CONFIGURED, got %v", resp.Error)
	}
}

// --- Password Reset Handler Tests ---

func TestForgotPassword_MissingEmail(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"email":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ForgotPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_EMAIL" {
		t.Errorf("expected MISSING_EMAIL, got %v", resp.Error)
	}
}

func TestForgotPassword_InvalidEmail(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"email":"not-an-email"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ForgotPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "INVALID_EMAIL" {
		t.Errorf("expected INVALID_EMAIL, got %v", resp.Error)
	}
}

func TestForgotPassword_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{broken`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", body)
	rec := httptest.NewRecorder()

	h.ForgotPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestResetPassword_MissingToken(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"token":"","password":"newpassword123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ResetPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_TOKEN" {
		t.Errorf("expected MISSING_TOKEN, got %v", resp.Error)
	}
}

func TestResetPassword_WeakPassword(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"token":"valid-token","password":"short"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ResetPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "WEAK_PASSWORD" {
		t.Errorf("expected WEAK_PASSWORD, got %v", resp.Error)
	}
}

func TestResetPassword_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{broken`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", body)
	rec := httptest.NewRecorder()

	h.ResetPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestResetPassword_PasswordTooLong(t *testing.T) {
	h := &Handlers{}

	longPassword := ""
	for i := 0; i < 200; i++ {
		longPassword += "a"
	}
	bodyStr := `{"token":"valid-token","password":"` + longPassword + `"}`
	body := bytes.NewBufferString(bodyStr)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ResetPassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "PASSWORD_TOO_LONG" {
		t.Errorf("expected PASSWORD_TOO_LONG, got %v", resp.Error)
	}
}

// --- Email Verification Handler Tests ---

func TestVerifyEmail_MissingToken(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"token":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-email", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.VerifyEmail(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_TOKEN" {
		t.Errorf("expected MISSING_TOKEN, got %v", resp.Error)
	}
}

func TestVerifyEmail_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{broken`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-email", body)
	rec := httptest.NewRecorder()

	h.VerifyEmail(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

// --- Change Password Handler Tests ---

func TestChangePassword_MissingCurrentPassword(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"current_password":"","new_password":"newpass123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ChangePassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_CURRENT_PASSWORD" {
		t.Errorf("expected MISSING_CURRENT_PASSWORD, got %v", resp.Error)
	}
}

func TestChangePassword_WeakNewPassword(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"current_password":"oldpass123","new_password":"short"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ChangePassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "WEAK_PASSWORD" {
		t.Errorf("expected WEAK_PASSWORD, got %v", resp.Error)
	}
}

func TestChangePassword_SamePassword(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"current_password":"samepass123","new_password":"samepass123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ChangePassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "SAME_PASSWORD" {
		t.Errorf("expected SAME_PASSWORD, got %v", resp.Error)
	}
}

func TestChangePassword_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{broken`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password", body)
	rec := httptest.NewRecorder()

	h.ChangePassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestChangePassword_NewPasswordTooLong(t *testing.T) {
	h := &Handlers{}

	longPassword := ""
	for i := 0; i < 200; i++ {
		longPassword += "x"
	}
	bodyStr := `{"current_password":"oldpass123","new_password":"` + longPassword + `"}`
	body := bytes.NewBufferString(bodyStr)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/change-password", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ChangePassword(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "PASSWORD_TOO_LONG" {
		t.Errorf("expected PASSWORD_TOO_LONG, got %v", resp.Error)
	}
}

// --- Delete Account Handler Tests ---

func TestDeleteAccount_MissingPassword(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"password":""}`)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/account", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.DeleteAccount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "MISSING_PASSWORD" {
		t.Errorf("expected MISSING_PASSWORD, got %v", resp.Error)
	}
}

func TestDeleteAccount_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{broken`)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/account", body)
	rec := httptest.NewRecorder()

	h.DeleteAccount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}
