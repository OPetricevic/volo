package intent

import (
	"encoding/json"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ModelParser uses a fine-tuned DistilBERT ONNX model for intent classification.
// It's the fallback when the rule-based parser has low confidence.
//
// NOTE: ONNX Runtime integration is platform-specific.
// On platforms without the runtime, the model parser is disabled gracefully.
// The actual ONNX inference is implemented in model_parser_onnx.go (Linux only).
type ModelParser struct {
	vocab    map[string]int64
	labels   []string
	maxLen   int
	padID    int64
	unkID    int64
	clsID    int64
	sepID    int64
	loaded   bool
	modelDir string
}

// ModelConfig matches the model-config.json exported by the training pipeline.
type ModelConfig struct {
	Labels   []string          `json:"labels"`
	Label2ID map[string]int    `json:"label2id"`
	ID2Label map[string]string `json:"id2label"`
	MaxLen   int               `json:"max_length"`
}

// NewModelParser loads the model configuration and vocabulary.
// Actual ONNX session creation is platform-dependent.
func NewModelParser(modelDir string) *ModelParser {
	mp := &ModelParser{
		maxLen:   64,
		padID:    0,
		unkID:    100,
		clsID:    101,
		sepID:    102,
		loaded:   false,
		modelDir: modelDir,
	}

	// Check if model exists
	onnxPath := filepath.Join(modelDir, "intent.onnx")
	if _, err := os.Stat(onnxPath); os.IsNotExist(err) {
		slog.Info("ONNX model not found, ML fallback disabled", "path", onnxPath)
		return mp
	}

	// Load config
	configPath := filepath.Join(modelDir, "model-config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		slog.Warn("model-config.json not found, ML fallback disabled", "error", err)
		return mp
	}

	var config ModelConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		slog.Warn("failed to parse model config", "error", err)
		return mp
	}
	mp.labels = config.Labels
	if config.MaxLen > 0 {
		mp.maxLen = config.MaxLen
	}

	// Load vocabulary
	vocabPath := filepath.Join(modelDir, "tokenizer", "vocab.txt")
	if err := mp.loadVocab(vocabPath); err != nil {
		slog.Warn("failed to load vocab", "error", err)
		return mp
	}

	// Try to initialize ONNX runtime (platform-dependent)
	if initOnnxSession(mp, onnxPath) {
		mp.loaded = true
		slog.Info("ONNX model loaded successfully", "labels", mp.labels, "maxLen", mp.maxLen)
	} else {
		slog.Info("ONNX runtime not available on this platform, ML fallback disabled")
	}

	return mp
}

// IsLoaded returns whether the model was successfully loaded.
func (mp *ModelParser) IsLoaded() bool {
	return mp.loaded
}

// Parse runs inference on the transcript and returns the predicted intent.
func (mp *ModelParser) Parse(transcript string) *Result {
	if !mp.loaded {
		return nil
	}

	// Tokenize
	inputIDs, attentionMask := mp.tokenize(transcript)

	// Run inference (platform-dependent)
	logits := runOnnxInference(mp, inputIDs, attentionMask)
	if logits == nil {
		return nil
	}

	// Softmax + argmax
	probs := softmaxF32(logits)
	maxIdx := 0
	maxProb := probs[0]
	for i := 1; i < len(probs); i++ {
		if probs[i] > maxProb {
			maxProb = probs[i]
			maxIdx = i
		}
	}

	if maxIdx >= len(mp.labels) {
		return nil
	}

	return &Result{
		Action:     mp.labels[maxIdx],
		Confidence: float64(maxProb),
	}
}

// tokenize converts text to token IDs using simplified WordPiece.
func (mp *ModelParser) tokenize(text string) ([]int64, []int64) {
	text = strings.ToLower(strings.TrimSpace(text))
	words := strings.Fields(text)

	ids := make([]int64, mp.maxLen)
	mask := make([]int64, mp.maxLen)

	ids[0] = mp.clsID
	mask[0] = 1
	pos := 1

	for _, word := range words {
		if pos >= mp.maxLen-1 {
			break
		}
		if id, ok := mp.vocab[word]; ok {
			ids[pos] = id
		} else {
			ids[pos] = mp.unkID
		}
		mask[pos] = 1
		pos++
	}

	if pos < mp.maxLen {
		ids[pos] = mp.sepID
		mask[pos] = 1
	}

	return ids, mask
}

func (mp *ModelParser) loadVocab(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	mp.vocab = make(map[string]int64)
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		token := strings.TrimSpace(line)
		if token != "" {
			mp.vocab[token] = int64(i)
		}
	}
	return nil
}

func softmaxF32(logits []float32) []float32 {
	max := logits[0]
	for _, v := range logits[1:] {
		if v > max {
			max = v
		}
	}
	probs := make([]float32, len(logits))
	sum := float32(0)
	for i, v := range logits {
		probs[i] = float32(math.Exp(float64(v - max)))
		sum += probs[i]
	}
	for i := range probs {
		probs[i] /= sum
	}
	return probs
}
