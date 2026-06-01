package middleware

import (
	"strings"
	"testing"
)

func TestValidateTranscript_Empty(t *testing.T) {
	err := ValidateTranscript("")
	if err == nil {
		t.Error("expected error for empty transcript")
	}
	if err.Code != "EMPTY_TRANSCRIPT" {
		t.Errorf("code = %q, want EMPTY_TRANSCRIPT", err.Code)
	}
}

func TestValidateTranscript_WhitespaceOnly(t *testing.T) {
	err := ValidateTranscript("   ")
	if err == nil {
		t.Error("expected error for whitespace-only transcript")
	}
}

func TestValidateTranscript_TooLong(t *testing.T) {
	long := strings.Repeat("a", 501)
	err := ValidateTranscript(long)
	if err == nil {
		t.Error("expected error for too-long transcript")
	}
	if err.Code != "TRANSCRIPT_TOO_LONG" {
		t.Errorf("code = %q, want TRANSCRIPT_TOO_LONG", err.Code)
	}
}

func TestValidateTranscript_Valid(t *testing.T) {
	err := ValidateTranscript("open youtube play lofi")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateTranscript_MaxLength(t *testing.T) {
	exact := strings.Repeat("a", 500)
	err := ValidateTranscript(exact)
	if err != nil {
		t.Error("expected no error for exactly max length")
	}
}

func TestValidateChatMessage_Empty(t *testing.T) {
	err := ValidateChatMessage("")
	if err == nil || err.Code != "EMPTY_MESSAGE" {
		t.Error("expected EMPTY_MESSAGE error")
	}
}

func TestValidateChatMessage_TooLong(t *testing.T) {
	long := strings.Repeat("x", 2001)
	err := ValidateChatMessage(long)
	if err == nil || err.Code != "MESSAGE_TOO_LONG" {
		t.Error("expected MESSAGE_TOO_LONG error")
	}
}

func TestValidateChatMessage_Valid(t *testing.T) {
	err := ValidateChatMessage("What did I search yesterday?")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateEmail_Empty(t *testing.T) {
	err := ValidateEmail("")
	if err == nil || err.Code != "MISSING_EMAIL" {
		t.Error("expected MISSING_EMAIL error")
	}
}

func TestValidateEmail_NoAt(t *testing.T) {
	err := ValidateEmail("notanemail")
	if err == nil || err.Code != "INVALID_EMAIL" {
		t.Error("expected INVALID_EMAIL error")
	}
}

func TestValidateEmail_NoDot(t *testing.T) {
	err := ValidateEmail("user@localhost")
	if err == nil || err.Code != "INVALID_EMAIL" {
		t.Error("expected INVALID_EMAIL error")
	}
}

func TestValidateEmail_Valid(t *testing.T) {
	err := ValidateEmail("user@example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateEmail_TooLong(t *testing.T) {
	long := strings.Repeat("a", 250) + "@b.com"
	err := ValidateEmail(long)
	if err == nil || err.Code != "EMAIL_TOO_LONG" {
		t.Error("expected EMAIL_TOO_LONG error")
	}
}

func TestValidatePassword_TooShort(t *testing.T) {
	err := ValidatePassword("short")
	if err == nil || err.Code != "WEAK_PASSWORD" {
		t.Error("expected WEAK_PASSWORD error")
	}
}

func TestValidatePassword_TooLong(t *testing.T) {
	long := strings.Repeat("a", 129)
	err := ValidatePassword(long)
	if err == nil || err.Code != "PASSWORD_TOO_LONG" {
		t.Error("expected PASSWORD_TOO_LONG error")
	}
}

func TestValidatePassword_Valid(t *testing.T) {
	err := ValidatePassword("securepassword123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidatePassword_ExactMin(t *testing.T) {
	err := ValidatePassword("12345678") // exactly 8
	if err != nil {
		t.Error("expected no error for exactly min length")
	}
}
