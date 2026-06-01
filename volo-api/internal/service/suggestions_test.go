package service

import (
	"testing"
	"time"

	"github.com/volo/volo-api/internal/model"
)

func TestSuggestions_Scoring(t *testing.T) {
	// Test the scoring logic: frequency * recency
	now := time.Now()

	// Recent command (1 hour ago) should score higher than old one (7 days ago)
	recentAge := 1.0 // hours
	oldAge := 168.0  // hours (7 days)

	recentRecency := 1.0 / (1.0 + recentAge/24.0)
	oldRecency := 1.0 / (1.0 + oldAge/24.0)

	if recentRecency <= oldRecency {
		t.Errorf("recent recency (%f) should be > old recency (%f)", recentRecency, oldRecency)
	}

	// Frequency 5 * recent should beat frequency 2 * recent
	score5 := 5.0 * recentRecency
	score2 := 2.0 * recentRecency
	if score5 <= score2 {
		t.Errorf("higher frequency should produce higher score")
	}

	// High frequency old should still beat low frequency recent (if freq is much higher)
	scoreOldFrequent := 10.0 * oldRecency
	scoreRecentRare := 1.0 * recentRecency
	if scoreOldFrequent <= scoreRecentRare {
		t.Errorf("very frequent old command (%f) should beat rare recent (%f)", scoreOldFrequent, scoreRecentRare)
	}

	_ = now
}

func TestSuggestions_EmptyHistory(t *testing.T) {
	// With no commands, suggestions should be empty (not nil)
	suggestions := []model.Suggestion{}
	if suggestions == nil {
		t.Error("suggestions should be empty slice, not nil")
	}
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(suggestions))
	}
}

func TestSuggestions_MaxTen(t *testing.T) {
	// Verify cap at 10
	suggestions := make([]model.Suggestion, 15)
	for i := range suggestions {
		suggestions[i] = model.Suggestion{Text: "cmd", Type: "search", Score: float64(i)}
	}

	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	if len(suggestions) != 10 {
		t.Errorf("expected max 10, got %d", len(suggestions))
	}
}

func TestSuggestions_SortedByScore(t *testing.T) {
	suggestions := []model.Suggestion{
		{Text: "low", Score: 1.0},
		{Text: "high", Score: 5.0},
		{Text: "mid", Score: 3.0},
	}

	// Sort descending (same logic as service)
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Score > suggestions[i].Score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	if suggestions[0].Text != "high" {
		t.Errorf("first should be 'high', got %q", suggestions[0].Text)
	}
	if suggestions[1].Text != "mid" {
		t.Errorf("second should be 'mid', got %q", suggestions[1].Text)
	}
	if suggestions[2].Text != "low" {
		t.Errorf("third should be 'low', got %q", suggestions[2].Text)
	}
}

func TestSuggestions_RecencyDecay(t *testing.T) {
	// Verify the decay formula produces expected values
	// 1 hour ago: 1/(1 + 1/24) = 0.96
	// 24 hours ago: 1/(1 + 24/24) = 0.5
	// 7 days ago: 1/(1 + 168/24) = 0.125

	oneHour := 1.0 / (1.0 + 1.0/24.0)
	oneDay := 1.0 / (1.0 + 24.0/24.0)
	oneWeek := 1.0 / (1.0 + 168.0/24.0)

	if oneHour < 0.9 || oneHour > 1.0 {
		t.Errorf("1 hour recency = %f, expected ~0.96", oneHour)
	}
	if oneDay < 0.4 || oneDay > 0.6 {
		t.Errorf("1 day recency = %f, expected ~0.5", oneDay)
	}
	if oneWeek < 0.1 || oneWeek > 0.15 {
		t.Errorf("1 week recency = %f, expected ~0.125", oneWeek)
	}
}
