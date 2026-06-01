package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/volo/volo-api/internal/model"
)

func TestProcessCommand_EmptyTranscript(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`{"transcript":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/command", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ProcessCommand(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}

	var resp model.Response
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Error == nil || resp.Error.Code != "EMPTY_TRANSCRIPT" {
		t.Errorf("expected EMPTY_TRANSCRIPT error, got %v", resp.Error)
	}
}

func TestProcessCommand_InvalidJSON(t *testing.T) {
	h := &Handlers{}

	body := bytes.NewBufferString(`not json`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/command", body)
	rec := httptest.NewRecorder()

	h.ProcessCommand(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestGetHistory_NoAuth(t *testing.T) {
	t.Skip("requires wired services — covered by integration tests")
}

func TestClearHistory_NoAuth(t *testing.T) {
	t.Skip("requires wired services — covered by integration tests")
}

func TestGetHistoryContext_NoAuth(t *testing.T) {
	t.Skip("requires wired services — covered by integration tests")
}
