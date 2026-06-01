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
