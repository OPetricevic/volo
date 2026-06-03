package service

import (
	"testing"
)

func TestHashToken_Deterministic(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"

	hash1 := hashToken(token)
	hash2 := hashToken(token)

	if hash1 != hash2 {
		t.Error("expected same hash for same token")
	}
}

func TestHashToken_DifferentTokens(t *testing.T) {
	hash1 := hashToken("token-a")
	hash2 := hashToken("token-b")

	if hash1 == hash2 {
		t.Error("expected different hashes for different tokens")
	}
}

func TestHashToken_Length(t *testing.T) {
	hash := hashToken("any-token")

	// SHA256 hex = 64 characters
	if len(hash) != 64 {
		t.Errorf("expected hash length 64, got %d", len(hash))
	}
}

func TestHashToken_NotEmpty(t *testing.T) {
	hash := hashToken("")
	// Even empty string should produce a valid hash
	if len(hash) != 64 {
		t.Errorf("expected hash length 64 for empty input, got %d", len(hash))
	}
}

// --- Password Reset Service Tests ---

func TestForgotPassword_TokenHashIs64Chars(t *testing.T) {
	// Verify that any raw token produces a 64-char SHA-256 hex hash
	rawToken := "dGVzdC10b2tlbi12YWx1ZS1mb3ItcmVzZXQ="
	hash := hashToken(rawToken)
	if len(hash) != 64 {
		t.Errorf("token hash length = %d, want 64", len(hash))
	}
}

func TestForgotPassword_TokenHashDeterministic(t *testing.T) {
	// Same raw token always produces the same hash (lookup consistency)
	raw := "base64-encoded-reset-token-value"
	hash1 := hashToken(raw)
	hash2 := hashToken(raw)
	if hash1 != hash2 {
		t.Error("reset token hash not deterministic")
	}
}

func TestForgotPassword_DifferentTokensDifferentHashes(t *testing.T) {
	hash1 := hashToken("token-for-user-a")
	hash2 := hashToken("token-for-user-b")
	if hash1 == hash2 {
		t.Error("different reset tokens should not produce same hash")
	}
}

func TestForgotPassword_EmptyTokenStillHashes(t *testing.T) {
	// Even an empty string should hash to a valid 64-char hex
	// (edge case protection — should never happen in practice)
	hash := hashToken("")
	if len(hash) != 64 {
		t.Errorf("empty token hash length = %d, want 64", len(hash))
	}
}
