# ML Model Training — Volo Intent Classifier

> How the DistilBERT intent classification model was trained, exported, and integrated into the Go API.

---

## Overview

Volo uses a **hybrid intent parser**:
1. **Rule-based parser** (~1ms) — handles 80%+ of commands via keyword pattern matching
2. **DistilBERT ONNX model** (~8ms) — fallback for ambiguous/novel commands when rule-based confidence < 0.80

The model classifies voice transcripts into 5 intents:

| Intent | Example |
|--------|---------|
| `search` | "search for react hooks", "what is kubernetes" |
| `navigate` | "open github", "go to youtube" |
| `open-and-search` | "search youtube for tutorials", "find recipes on reddit" |
| `open-and-play` | "open youtube play lofi", "play jazz on spotify" |
| `browser-control` | "close tab", "scroll down", "go back" |

---

## Architecture

```
DistilBERT (67M params, distilbert-base-uncased)
        ↓ Fine-tuned (sequence classification, 5 labels)
PyTorch model
        ↓ Exported via HuggingFace Optimum
ONNX model (255 MB)
        ↓ Dynamic quantization (INT8)
Quantized ONNX (64 MB) → loaded in Go via onnxruntime-go
```

---

## Training Pipeline

### Prerequisites

- Docker running locally
- `docker-compose.dev.yml` at project root

### Step 1: Start the ML container

```bash
docker compose -f docker-compose.dev.yml up -d model
```

This builds a Python 3.11 container with PyTorch (CPU-only), HuggingFace Transformers, ONNX Runtime, and JupyterLab.

### Step 2: Prepare dataset

```bash
docker compose -f docker-compose.dev.yml exec model python scripts/prepare_dataset.py
```

**What it does:**
- Attempts to load the SNIPS NLU dataset from HuggingFace Hub (graceful fallback if unavailable)
- Generates 370+ custom browser command examples from templates
- Augments with filler words ("um", "uh", "like") to simulate real speech recognition output
- Balances classes (max 400 per intent)
- Splits 85/15 into train/eval

**Output:**
```
data/train.json  — 377 training examples
data/eval.json   — 67 evaluation examples
```

### Step 3: Train the model

```bash
docker compose -f docker-compose.dev.yml exec model python scripts/train.py
```

**Configuration:**
- Base model: `distilbert-base-uncased`
- Epochs: 8
- Batch size: 32
- Learning rate: 3e-5 (linear warmup 10%)
- Max sequence length: 64 tokens
- Training time: ~5 minutes on CPU

**Results:**
```
Epoch 5: accuracy=1.0, f1_macro=1.0
Epoch 8: accuracy=1.0, f1_macro=1.0 (converged)

Final evaluation:
  search:          precision=1.00, recall=1.00, f1=1.00 (32 samples)
  navigate:        precision=1.00, recall=1.00, f1=1.00 (9 samples)
  open-and-search: precision=1.00, recall=1.00, f1=1.00 (12 samples)
  open-and-play:   precision=1.00, recall=1.00, f1=1.00 (12 samples)
  browser-control: precision=1.00, recall=1.00, f1=1.00 (2 samples)
```

> Note: 100% eval accuracy is expected with a small, clean dataset where patterns are distinct.
> Real-world accuracy depends on speech recognition noise, which the augmentation partially simulates.

### Step 4: Export to ONNX

```bash
docker compose -f docker-compose.dev.yml exec model python scripts/export_onnx.py
```

**What it does:**
1. Exports PyTorch model → ONNX format via HuggingFace Optimum
2. Quantizes to INT8 (dynamic quantization, AVX-512)
3. Saves tokenizer config for Go integration
4. Verifies the exported model with a test inference

**Output:**
```
models/volo-intent.onnx           — 255.5 MB (full precision)
models/volo-intent-quantized.onnx — 64.3 MB (quantized, used in production)
models/tokenizer/                 — tokenizer.json, vocab.txt, config
models/model-config.json          — label mapping, vocab size, max_length
```

**Verification result:**
```
Input: "open youtube play lofi"
Predicted: open-and-play (confidence: 0.531)
```

### Step 5: Copy to API

```bash
# From project root
cp volo-model/models/volo-intent-quantized.onnx volo-api/models/intent.onnx
cp -r volo-model/models/tokenizer volo-api/models/tokenizer
cp volo-model/models/model-config.json volo-api/models/model-config.json
```

On Windows (PowerShell):
```powershell
Copy-Item "volo-model\models\volo-intent-quantized.onnx" "volo-api\models\intent.onnx"
Copy-Item "volo-model\models\model-config.json" "volo-api\models\model-config.json"
Copy-Item "volo-model\models\tokenizer\*" "volo-api\models\tokenizer\" -Recurse
```

### Step 6: Stop the container

```bash
docker compose -f docker-compose.dev.yml stop model
```

---

## Go Integration

The model is loaded in `volo-api/internal/intent/model_parser.go` via `onnxruntime-go`.

**Platform notes:**
- **Linux:** Full ONNX Runtime integration (production deployment)
- **Windows/macOS:** Uses `onnx_stub.go` — the model parser returns a zero-confidence result, letting the rule-based parser handle everything

To enable on Linux, replace the stub with the real implementation and set:
```bash
VOLO_MODEL_PATH=./models
```

---

## Issues Encountered During Training

| Issue | Cause | Resolution |
|-------|-------|------------|
| Docker build timeout downloading PyTorch | Default PyTorch install pulls CUDA libs (~900MB) | Used `--index-url https://download.pytorch.org/whl/cpu` for CPU-only build (192MB) |
| SNIPS dataset `HfUriError: Repository id must be 'namespace/name'` | HuggingFace `datasets` library v4.8+ requires namespaced dataset IDs | Made SNIPS loading graceful with try/except fallback — custom data is sufficient |
| `ModuleNotFoundError: No module named 'optimum.onnxruntime'` | The `optimum` package needs the `[onnxruntime]` extra | Ran `pip install optimum[onnxruntime]` inside the container |

---

## Retraining

Retrain when:
- Adding new intent labels
- Adding significantly more training data
- Improving accuracy on edge cases

Steps:
1. Edit `volo-model/data/train.json` (or modify `prepare_dataset.py`)
2. Re-run steps 3-5 above
3. Restart the Go API (or hot-swap the ONNX file on Linux)

---

## File Locations

```
volo-model/
├── scripts/
│   ├── prepare_dataset.py  — Dataset generation + SNIPS loading
│   ├── train.py            — DistilBERT fine-tuning
│   └── export_onnx.py     — ONNX export + quantization + verification
├── data/
│   ├── train.json          — 377 training examples
│   └── eval.json           — 67 evaluation examples
├── models/
│   ├── volo-intent/        — Full PyTorch checkpoint
│   ├── volo-intent.onnx    — Exported ONNX (255 MB)
│   ├── volo-intent-quantized.onnx — Quantized (64 MB)
│   ├── tokenizer/          — HuggingFace tokenizer files
│   └── model-config.json   — Label mapping + config
├── Dockerfile              — Python 3.11 + PyTorch CPU + JupyterLab
└── requirements.txt        — ML dependencies

volo-api/models/
├── intent.onnx             — Quantized model (copied from volo-model)
├── tokenizer/              — Tokenizer files (copied from volo-model)
└── model-config.json       — Config (copied from volo-model)
```

---

## References

- [DistilBERT paper](https://arxiv.org/abs/1910.01108) — Sanh et al., 2019
- [ONNX Runtime](https://onnxruntime.ai/) — Cross-platform inference engine
- [HuggingFace Optimum](https://huggingface.co/docs/optimum/) — ONNX export tooling
- [onnxruntime-go](https://github.com/yalue/onnxruntime_go) — Go bindings for ONNX Runtime
