package model

import (
	"encoding/json"
	"testing"
)

func TestResponse_SuccessJSON(t *testing.T) {
	resp := Response{
		Data: map[string]string{"status": "ok"},
	}

	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(b, &decoded)

	if decoded["error"] != nil {
		t.Error("expected no error field in success response")
	}
	data := decoded["data"].(map[string]interface{})
	if data["status"] != "ok" {
		t.Errorf("data.status = %v, want ok", data["status"])
	}
}

func TestResponse_ErrorJSON(t *testing.T) {
	resp := Response{
		Error: &ErrorResponse{
			Code:          "TEST_ERROR",
			Message:       "Something failed.",
			InternalError: "handler.X → service.Y: actual error",
		},
	}

	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded map[string]interface{}
	json.Unmarshal(b, &decoded)

	if decoded["data"] != nil {
		t.Error("expected no data field in error response")
	}

	errObj := decoded["error"].(map[string]interface{})
	if errObj["code"] != "TEST_ERROR" {
		t.Errorf("code = %v, want TEST_ERROR", errObj["code"])
	}
	if errObj["message"] != "Something failed." {
		t.Errorf("message = %v", errObj["message"])
	}
	if errObj["internal_error"] != "handler.X → service.Y: actual error" {
		t.Errorf("internal_error = %v", errObj["internal_error"])
	}
}

func TestUser_PasswordNeverSerialized(t *testing.T) {
	hash := "$2a$10$abc123"
	cred := Credential{
		ID:           "cred-1",
		UserID:       "user-1",
		Provider:     "password",
		PasswordHash: &hash,
	}

	b, _ := json.Marshal(cred)
	var decoded map[string]interface{}
	json.Unmarshal(b, &decoded)

	if _, exists := decoded["password_hash"]; exists {
		t.Error("password_hash should never appear in JSON output")
	}
}

func TestSession_TokenHashNeverSerialized(t *testing.T) {
	session := Session{
		ID:        "sess-1",
		UserID:    "user-1",
		TokenHash: "sha256-hash-value",
	}

	b, _ := json.Marshal(session)
	var decoded map[string]interface{}
	json.Unmarshal(b, &decoded)

	if _, exists := decoded["token_hash"]; exists {
		t.Error("token_hash should never appear in JSON output")
	}
}
