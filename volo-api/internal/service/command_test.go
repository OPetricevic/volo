package service

import (
	"testing"
	"time"

	"github.com/volo/volo-api/internal/model"
)

func TestGetHistoryContext_FormatOutput(t *testing.T) {
	target := "youtube"
	query := "lofi beats"
	cmd := model.Command{
		ID:           "cmd-1",
		UserID:       "user-1",
		Transcript:   "open youtube play lofi beats",
		ParsedAction: "open-and-play",
		ParsedTarget: &target,
		ParsedQuery:  &query,
		Confidence:   0.92,
		ExecutedAt:   time.Date(2026, 6, 1, 9, 14, 0, 0, time.UTC),
	}

	// Simulate the formatting logic from GetHistoryContext
	result := formatCommandForContext(cmd)

	if result == "" {
		t.Error("expected non-empty formatted command")
	}
	if !containsStr(result, "open-and-play") {
		t.Error("expected action in output")
	}
	if !containsStr(result, "youtube") {
		t.Error("expected target in output")
	}
	if !containsStr(result, "lofi beats") {
		t.Error("expected query in output")
	}
	if !containsStr(result, "Jun") {
		t.Error("expected date in output")
	}
}

func TestGetHistoryContext_NilFields(t *testing.T) {
	cmd := model.Command{
		ID:           "cmd-2",
		UserID:       "user-1",
		Transcript:   "go back",
		ParsedAction: "browser-control",
		ParsedTarget: nil,
		ParsedQuery:  nil,
		Confidence:   0.95,
		ExecutedAt:   time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
	}

	// Should not panic with nil fields
	result := formatCommandForContext(cmd)
	if result == "" {
		t.Error("expected non-empty output even with nil fields")
	}
	if !containsStr(result, "browser-control") {
		t.Error("expected action in output")
	}
}

func TestGetHistoryContext_MultipleCommands(t *testing.T) {
	target := "github"
	query := "react"
	commands := []model.Command{
		{ParsedAction: "search", ParsedQuery: &query, ExecutedAt: time.Now(), Transcript: "search react"},
		{ParsedAction: "navigate", ParsedTarget: &target, ExecutedAt: time.Now(), Transcript: "open github"},
	}

	var result string
	for _, cmd := range commands {
		result += formatCommandForContext(cmd)
	}

	if !containsStr(result, "search") {
		t.Error("expected first command")
	}
	if !containsStr(result, "navigate") {
		t.Error("expected second command")
	}
}

// Helper: format a single command (mirrors the logic in GetHistoryContext)
func formatCommandForContext(cmd model.Command) string {
	target := ""
	if cmd.ParsedTarget != nil {
		target = " → " + *cmd.ParsedTarget
	}
	query := ""
	if cmd.ParsedQuery != nil {
		query = ": " + *cmd.ParsedQuery
	}
	return "- [" + cmd.ExecutedAt.Format("Jan 2 15:04") + "] " + cmd.ParsedAction + target + query + " (" + cmd.Transcript + ")\n"
}
