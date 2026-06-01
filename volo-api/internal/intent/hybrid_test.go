package intent

import (
	"testing"
)

func TestHybridParser_RuleBasedOnly(t *testing.T) {
	// With no model directory, hybrid parser uses rule-based only
	h := NewHybridParser("")

	if h.IsModelLoaded() {
		t.Error("expected model not loaded with empty dir")
	}

	// Should still parse correctly using rules
	result := h.Parse("open youtube")
	if result.Action != "navigate" {
		t.Errorf("action = %q, want navigate", result.Action)
	}
	if result.Target == nil || *result.Target != "youtube" {
		t.Error("expected target youtube")
	}
}

func TestHybridParser_HighConfidence_NoFallback(t *testing.T) {
	h := NewHybridParser("") // no model

	// High confidence commands should not need the model
	result := h.Parse("go back")
	if result.Action != "browser-control" {
		t.Errorf("action = %q, want browser-control", result.Action)
	}
	if result.Confidence < 0.80 {
		t.Errorf("confidence = %f, expected >= 0.80", result.Confidence)
	}
}

func TestHybridParser_LowConfidence_WouldFallback(t *testing.T) {
	h := NewHybridParser("") // no model, so fallback returns rule-based result anyway

	// Ambiguous command — rule-based gives low confidence
	result := h.Parse("play some chill music")
	if result.Action == "" {
		t.Error("expected non-empty action")
	}
	// Without model, we get the rule-based fallback (search with low confidence)
	if result.Confidence >= 0.80 {
		t.Errorf("expected low confidence for ambiguous input, got %f", result.Confidence)
	}
}

func TestHybridParser_NonexistentModelDir(t *testing.T) {
	// Should not panic with a nonexistent directory
	h := NewHybridParser("/nonexistent/path/to/models")

	if h.IsModelLoaded() {
		t.Error("expected model not loaded with nonexistent dir")
	}

	// Should still work with rule-based
	result := h.Parse("search react hooks")
	if result.Action != "search" {
		t.Errorf("action = %q, want search", result.Action)
	}
}

func TestHybridParser_ConfidenceThreshold(t *testing.T) {
	// Verify the threshold constant
	if confidenceThreshold != 0.80 {
		t.Errorf("threshold = %f, want 0.80", confidenceThreshold)
	}
}
