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
