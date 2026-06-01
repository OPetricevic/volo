package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/volo/volo-api/internal/model"
)

func TestHealth(t *testing.T) {
	h := &Handlers{}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	var resp model.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}
	if data["status"] != "ok" {
		t.Errorf("status = %v, want ok", data["status"])
	}
}

func TestRespondError(t *testing.T) {
	rec := httptest.NewRecorder()
	respondError(rec, http.StatusBadRequest, "TEST_ERROR", "Something went wrong.", "internal details here")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}
	if resp.Error.Code != "TEST_ERROR" {
		t.Errorf("code = %q, want TEST_ERROR", resp.Error.Code)
	}
	if resp.Error.Message != "Something went wrong." {
		t.Errorf("message = %q, want 'Something went wrong.'", resp.Error.Message)
	}
	if resp.Error.InternalError != "internal details here" {
		t.Errorf("internal_error = %q, want 'internal details here'", resp.Error.InternalError)
	}
}

func TestDecode(t *testing.T) {
	body := bytes.NewBufferString(`{"transcript":"hello world"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")

	var cmd model.ProcessCommandRequest
	if err := decode(req, &cmd); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if cmd.Transcript != "hello world" {
		t.Errorf("transcript = %q, want 'hello world'", cmd.Transcript)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)

	var cmd model.ProcessCommandRequest
	if err := decode(req, &cmd); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
