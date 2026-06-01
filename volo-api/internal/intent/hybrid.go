package intent

import "log/slog"

const confidenceThreshold = 0.80

// HybridParser combines rule-based (fast) and ML model (fallback) parsing.
type HybridParser struct {
	rules *Parser
	model *ModelParser
}

// NewHybridParser creates a hybrid parser.
// If modelDir is empty or model fails to load, only rule-based parsing is used.
func NewHybridParser(modelDir string) *HybridParser {
	rules := NewParser()

	var model *ModelParser
	if modelDir != "" {
		model = NewModelParser(modelDir)
		if model.IsLoaded() {
			slog.Info("hybrid parser: rule-based + ML model active")
		} else {
			slog.Info("hybrid parser: ML model not available, rule-based only")
			model = nil
		}
	} else {
		slog.Info("hybrid parser: no model directory configured, rule-based only")
	}

	return &HybridParser{rules: rules, model: model}
}

// Parse tries rule-based first. If confidence is below threshold, falls back to ML model.
func (h *HybridParser) Parse(transcript string) Result {
	// 1. Always try rule-based first (fast path, ~1ms)
	result := h.rules.Parse(transcript)

	// 2. If confident enough, use it
	if result.Confidence >= confidenceThreshold {
		return result
	}

	// 3. Fall back to ML model (if available)
	if h.model != nil {
		mlResult := h.model.Parse(transcript)
		if mlResult != nil && mlResult.Confidence > result.Confidence {
			slog.Debug("ML model override",
				"transcript", transcript,
				"rule_action", result.Action,
				"rule_conf", result.Confidence,
				"ml_action", mlResult.Action,
				"ml_conf", mlResult.Confidence,
			)
			// ML model gives us the intent but not the URL — use rule-based for that
			// or return as-is (client can build URL from action + target)
			return *mlResult
		}
	}

	// 4. Return whatever we have (even low confidence)
	return result
}

// IsModelLoaded returns whether the ML model is active.
func (h *HybridParser) IsModelLoaded() bool {
	return h.model != nil && h.model.IsLoaded()
}
