package service

import "testing"

func TestEmailService_NotConfigured(t *testing.T) {
	svc := NewEmailService("")
	if svc.IsConfigured() {
		t.Error("expected service to be not configured with empty API key")
	}
}

func TestEmailService_Configured(t *testing.T) {
	svc := NewEmailService("re_test_key_123")
	if !svc.IsConfigured() {
		t.Error("expected service to be configured with API key")
	}
}

func TestEmailService_SendWithoutKey(t *testing.T) {
	// Should not error — just logs a warning and returns nil
	svc := NewEmailService("")
	err := svc.SendPasswordReset("user@test.com", "token123", "https://volo.app/reset")
	if err != nil {
		t.Errorf("expected no error when unconfigured, got: %v", err)
	}
}
