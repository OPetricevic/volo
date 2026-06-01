package service

import (
	"testing"
	"time"

	"github.com/volo/volo-api/internal/model"
)

func TestBuildHistoryContext_Empty(t *testing.T) {
	result := buildHistoryContext(nil)
	if result != "(No commands in history yet)" {
		t.Errorf("expected empty history message, got: %q", result)
	}

	result2 := buildHistoryContext([]model.Command{})
	if result2 != "(No commands in history yet)" {
		t.Errorf("expected empty history message, got: %q", result2)
	}
}

func TestBuildHistoryContext_WithCommands(t *testing.T) {
	target := "youtube"
	query := "lofi beats"
	commands := []model.Command{
		{
			ID:           "cmd-1",
			UserID:       "user-1",
			Transcript:   "open youtube play lofi beats",
			ParsedAction: "open-and-play",
			ParsedTarget: &target,
			ParsedQuery:  &query,
			Confidence:   0.92,
			ExecutedAt:   time.Date(2026, 6, 1, 9, 14, 0, 0, time.UTC),
		},
		{
			ID:           "cmd-2",
			UserID:       "user-1",
			Transcript:   "search react hooks",
			ParsedAction: "search",
			ParsedQuery:  &query,
			Confidence:   0.9,
			ExecutedAt:   time.Date(2026, 6, 1, 10, 30, 0, 0, time.UTC),
		},
	}

	result := buildHistoryContext(commands)

	if result == "(No commands in history yet)" {
		t.Error("expected history context, got empty message")
	}
	if len(result) == 0 {
		t.Error("expected non-empty context string")
	}
	// Should contain the action types
	if !containsStr(result, "open-and-play") {
		t.Error("expected context to contain 'open-and-play'")
	}
	if !containsStr(result, "search") {
		t.Error("expected context to contain 'search'")
	}
	// Should contain the target
	if !containsStr(result, "youtube") {
		t.Error("expected context to contain 'youtube'")
	}
}

func TestBuildHistoryContext_NilFields(t *testing.T) {
	// Commands with nil target and query should not panic
	commands := []model.Command{
		{
			ID:           "cmd-1",
			UserID:       "user-1",
			Transcript:   "go back",
			ParsedAction: "browser-control",
			ParsedTarget: nil,
			ParsedQuery:  nil,
			Confidence:   0.95,
			ExecutedAt:   time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC),
		},
	}

	// Should not panic
	result := buildHistoryContext(commands)
	if len(result) == 0 {
		t.Error("expected non-empty context")
	}
	if !containsStr(result, "browser-control") {
		t.Error("expected context to contain 'browser-control'")
	}
}

func TestBuildChatPrompt(t *testing.T) {
	history := "- [Jun 1 09:14] open-and-play → youtube: lofi beats"
	message := "What did I listen to today?"

	prompt := buildChatPrompt(message, history)

	// Should contain the system instruction
	if !containsStr(prompt, "You are Volo") {
		t.Error("expected prompt to contain system instruction")
	}
	// Should contain the history
	if !containsStr(prompt, "lofi beats") {
		t.Error("expected prompt to contain history context")
	}
	// Should contain the user message
	if !containsStr(prompt, "What did I listen to today?") {
		t.Error("expected prompt to contain user message")
	}
	// Should end with the assistant prompt marker
	if !containsStr(prompt, "Volo:") {
		t.Error("expected prompt to end with 'Volo:'")
	}
}

func TestBuildChatPrompt_EmptyHistory(t *testing.T) {
	prompt := buildChatPrompt("hello", "(No commands in history yet)")
	if !containsStr(prompt, "No commands in history yet") {
		t.Error("expected prompt to indicate empty history")
	}
}

func TestChatRequest_Validation(t *testing.T) {
	svc := &ChatService{repos: nil}

	// Empty message should fail
	_, err := svc.Process(nil, "user-1", ChatRequest{Message: ""})
	if err == nil {
		t.Error("expected error for empty message")
	}
}

func TestGenerateConversationID(t *testing.T) {
	id1 := GenerateConversationID()
	id2 := GenerateConversationID()

	if id1 == "" {
		t.Error("expected non-empty conversation ID")
	}
	if id1 == id2 {
		t.Error("expected unique conversation IDs")
	}
	// Should be UUID format (36 chars with hyphens)
	if len(id1) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(id1))
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
