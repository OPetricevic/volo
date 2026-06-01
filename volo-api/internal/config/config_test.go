package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any env vars that might interfere
	os.Unsetenv("VOLO_PORT")
	os.Unsetenv("VOLO_DATABASE_URL")
	os.Unsetenv("VOLO_REDIS_ADDR")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("RedisAddr = %q, want localhost:6379", cfg.RedisAddr)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want info", cfg.LogLevel)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("VOLO_PORT", "9090")
	os.Setenv("VOLO_LOG_LEVEL", "debug")
	defer os.Unsetenv("VOLO_PORT")
	defer os.Unsetenv("VOLO_LOG_LEVEL")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want 9090", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", cfg.LogLevel)
	}
}

func TestSplitEnv_Wildcard(t *testing.T) {
	result := splitEnv("NONEXISTENT_KEY", "*")
	if len(result) != 1 || result[0] != "*" {
		t.Errorf("expected [*], got %v", result)
	}
}

func TestSplitEnv_MultipleOrigins(t *testing.T) {
	os.Setenv("TEST_ORIGINS", "https://volo.app,chrome-extension://abc,http://localhost:3000")
	defer os.Unsetenv("TEST_ORIGINS")

	result := splitEnv("TEST_ORIGINS", "*")
	if len(result) != 3 {
		t.Fatalf("expected 3 origins, got %d: %v", len(result), result)
	}
	if result[0] != "https://volo.app" {
		t.Errorf("result[0] = %q, want https://volo.app", result[0])
	}
	if result[1] != "chrome-extension://abc" {
		t.Errorf("result[1] = %q, want chrome-extension://abc", result[1])
	}
	if result[2] != "http://localhost:3000" {
		t.Errorf("result[2] = %q, want http://localhost:3000", result[2])
	}
}
