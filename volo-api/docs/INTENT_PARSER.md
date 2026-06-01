# Intent Parser

## Overview

The intent parser converts raw voice transcripts into structured commands. It uses a **hybrid approach**: a fast rule-based parser handles obvious commands, and a fine-tuned DistilBERT model (ONNX) handles ambiguous or novel inputs.

```
"open youtube play lofi hip hop"
            │
            ▼
┌───────────────────────────────────┐
│       Rule-Based Parser            │  ~1ms
│                                    │
│  Regex patterns + keyword matching │
│  + fuzzy site resolution           │
│                                    │
│  confidence = 0.92 ✓               │──→ Execute immediately
└───────────────────────────────────┘
            │
     confidence < 0.80
            │
            ▼
┌───────────────────────────────────┐
│       DistilBERT (ONNX)           │  ~8ms
│                                    │
│  Fine-tuned intent classifier      │
│  + slot extraction                 │
│                                    │
│  confidence = 0.88 ✓               │──→ Execute
└───────────────────────────────────┘
```

---

## Intent Categories

| Intent | Description | Example |
|--------|-------------|---------|
| `search` | Google search for a query | "search react hooks tutorial" |
| `navigate` | Open a known site | "open youtube", "go to github" |
| `open-and-search` | Open site + search within it | "open youtube lofi beats" |
| `open-and-play` | Open site + search + autoplay | "open youtube play lofi" |
| `browser-control` | Tab/navigation control | "go back", "close tab", "scroll down" |

---

## Slot Types

| Slot | Description | Example values |
|------|-------------|---------------|
| `target` | The site/app to interact with | youtube, github, spotify, gmail |
| `query` | Search terms or content to find | "lofi beats", "react hooks", "jazz" |

---

## Rule-Based Parser

### Location

- **Extension:** `src/shared/parser.ts` (TypeScript, runs client-side)
- **API:** `internal/intent/parser.go` (Go, runs server-side)

Both implementations are functionally identical — same rules, same priority order.

### Priority Order

The parser tries patterns in this order. First match wins.

```
1. Browser control (exact match)
   "go back" → { action: "browser-control", query: "back" }
   "close tab" → { action: "browser-control", query: "close-tab" }
   "scroll down" → { action: "browser-control", query: "scroll-down" }

2. Open + play (regex: open [site] play [query])
   "open youtube play lofi" → { action: "open-and-play", target: "youtube", query: "lofi" }

3. Open + search (regex: open [site] [query])
   "open youtube lofi beats" → { action: "open-and-search", target: "youtube", query: "lofi beats" }

4. Navigate (regex: open|go to|goto [site])
   "open youtube" → { action: "navigate", target: "youtube" }
   "open example.com" → { action: "navigate", target: "https://example.com" }

5. Search (regex: search for|search|look up|find [query])
   "search for react hooks" → { action: "search", query: "react hooks" }

6. Fallback (entire text becomes a Google search)
   "where does burek come from" → { action: "search", query: "where does burek come from", confidence: 0.6 }
```

### Site Resolution

When the parser encounters a site name, it resolves it through three layers:

```
Input: "tube"
  │
  ├─ 1. Direct match in SITE_MAP?     "tube" → ✗
  ├─ 2. Alias match?                   "tube" → "youtube" ✓
  └─ 3. Partial/contains match?        (fallback)

Aliases:
  tube → youtube
  yt   → youtube
  gh   → github
  tw   → twitter
  so   → stackoverflow
```

### URL Building

After parsing, the result includes an `execute_url` for the client to navigate to:

| Action | URL Pattern |
|--------|-------------|
| search | `https://www.google.com/search?q={query}` |
| navigate | `SITE_MAP[target]` (e.g., `https://www.youtube.com`) |
| open-and-search | `SEARCH_MAP[target](query)` (e.g., YouTube search URL) |
| open-and-play | Same as open-and-search (client handles autoplay) |
| browser-control | No URL — executed via Chrome APIs |

### Confidence Scores

| Match type | Confidence |
|-----------|-----------|
| Browser control (exact) | 0.95 |
| Open + play (regex + known site) | 0.92 |
| Navigate (known site) | 0.92 |
| Open + search (regex + known site) | 0.90 |
| Search (explicit keyword) | 0.90 |
| Navigate (raw URL with dot) | 0.85 |
| Fallback (unknown pattern) | 0.60 |

The threshold for falling back to the ML model is **0.80**. Anything below that gets sent to DistilBERT.

---

## DistilBERT Model (Phase 3)

### When It's Used

The model is only called when:
1. Rule-based confidence < 0.80 (ambiguous command)
2. No regex pattern matched (novel phrasing)

Examples that would trigger the model:
- "find that video about golang I watched yesterday"
- "play some chill music"
- "take me to my email"
- "what's the weather"

### Architecture

```
┌─────────────────────────────────────────────────┐
│  DistilBERT (67M params)                         │
│                                                  │
│  Input: tokenized transcript                     │
│                                                  │
│  ┌────────────────────────────────────────────┐  │
│  │  [CLS] play some chill music [SEP]         │  │
│  └────────────────────────────────────────────┘  │
│                    │                             │
│         ┌─────────┴─────────┐                    │
│         ▼                   ▼                    │
│  ┌─────────────┐    ┌─────────────┐             │
│  │ Intent Head │    │  Slot Head  │             │
│  │ (classify)  │    │ (per-token) │             │
│  └─────────────┘    └─────────────┘             │
│         │                   │                    │
│         ▼                   ▼                    │
│  "open-and-play"    target: null                 │
│  confidence: 0.88   query: "chill music"         │
└─────────────────────────────────────────────────┘
```

### Two-Head Design

1. **Intent classification head** — Softmax over 5 intent labels. Uses the `[CLS]` token embedding.
2. **Slot extraction head** — Per-token BIO tagging (B-target, I-target, B-query, I-query, O). Identifies which words are the target and which are the query.

### Training Data Requirements

~500-1000 labeled examples covering:
- Simple commands (same as rule-based handles)
- Ambiguous commands (the model's real value)
- Variations and paraphrases
- Edge cases (partial site names, slang, typos from speech recognition)

See [volo-model/docs/DATASET.md](../../volo-model/docs/DATASET.md) for format details.

### ONNX Integration in Go

```go
// Loaded once at startup
type ModelParser struct {
    session   *ort.AdvancedSession
    tokenizer *WordPieceTokenizer
    labels    []string  // ["search", "navigate", "open-and-search", ...]
}

// Called only when rule-based confidence < 0.80
func (m *ModelParser) Parse(transcript string) Result {
    // 1. Tokenize (WordPiece, same as HuggingFace)
    tokens := m.tokenizer.Encode(transcript, 64) // max 64 tokens

    // 2. Run ONNX inference
    intentLogits, slotLogits := m.session.Run(tokens)

    // 3. Decode intent (argmax of softmax)
    intent := m.labels[argmax(intentLogits)]
    confidence := softmax(intentLogits)[argmax(intentLogits)]

    // 4. Decode slots (BIO tags → spans)
    target, query := decodeSlots(slotLogits, tokens)

    return Result{
        Action:     intent,
        Target:     target,
        Query:      query,
        Confidence: confidence,
    }
}
```

### Performance

| Metric | Value |
|--------|-------|
| Model size (quantized) | ~65 MB |
| Inference time (CPU) | ~8ms |
| Memory usage | ~150 MB |
| Max input length | 64 tokens (~50 words) |

This is fast enough to run on every request without noticeable latency. The rule-based parser still handles 80%+ of commands at ~1ms.

---

## Hybrid Flow (Complete)

```go
func (p *HybridParser) Parse(transcript string) Result {
    // 1. Always try rule-based first (fast path)
    result := p.ruleBased.Parse(transcript)

    // 2. If confident enough, use it
    if result.Confidence >= 0.80 {
        return result
    }

    // 3. Fall back to ML model
    if p.model != nil {
        mlResult := p.model.Parse(transcript)
        // Use ML result if it's more confident
        if mlResult.Confidence > result.Confidence {
            return mlResult
        }
    }

    // 4. Return whatever we have (even low confidence)
    return result
}
```

---

## Continuous Improvement

The system gets better over time without retraining:

### Short-term (frequency model)
- Track which commands a user runs most often
- Boost those in suggestions
- Time-of-day weighting (morning = email, evening = youtube)

### Medium-term (correction logging)
- If user says something, gets a result, then immediately says it differently → log as correction
- These corrections become training data for the next model version

### Long-term (model retraining)
- Batch new training examples from corrections + new patterns
- Retrain model (10 min on CPU)
- Export new ONNX file
- Hot-swap on the server (`docker compose restart volo-api`)

---

## Testing the Parser

### Unit Tests (both implementations)

```bash
# Go (28 tests)
go test ./internal/intent/ -v

# TypeScript (41 tests including parser)
cd volo-extension && npm test
```

### What's Tested

- All 5 intent categories with multiple phrasings
- Fuzzy site matching (aliases, partial matches)
- Case insensitivity
- Raw URL detection
- Confidence score ranges
- URL building for each action type
- Edge cases (empty input, very long input, special characters)

### Manual Testing

For the ML model (once trained):
```bash
# Test inference directly
cd volo-model
python scripts/evaluate.py --input "play some jazz on spotify"
# → { intent: "open-and-play", target: "spotify", query: "jazz", confidence: 0.91 }
```
