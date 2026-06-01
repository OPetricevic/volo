# Volo Model

ML training pipeline for the Volo intent classifier. Fine-tunes DistilBERT for voice command classification, exports to ONNX for the Go backend.

## Quick Start (Docker — recommended)

From the project root:

```bash
docker compose -f docker-compose.dev.yml up -d model
```

Then open Jupyter Lab at **http://localhost:8888** (token: `volo`).

Or run scripts directly:

```bash
docker compose -f docker-compose.dev.yml exec model python scripts/prepare_dataset.py
docker compose -f docker-compose.dev.yml exec model python scripts/train.py
docker compose -f docker-compose.dev.yml exec model python scripts/export_onnx.py
```

## Workflow

```
1. Prepare dataset     →  data/train.json, data/eval.json
2. Train model         →  models/volo-intent/ (HuggingFace checkpoint)
3. Export to ONNX      →  models/volo-intent-quantized.onnx (~65MB)
4. Copy to API         →  ../volo-api/models/
```

### Step 1: Prepare Dataset

```bash
python scripts/prepare_dataset.py
```

Combines SNIPS dataset (remapped to Volo intents) with custom browser command examples. Outputs balanced train/eval splits.

### Step 2: Train

```bash
python scripts/train.py
```

Fine-tunes `distilbert-base-uncased` for ~8 epochs. Takes ~10 minutes on CPU, ~2 minutes on GPU. Outputs accuracy and F1 scores.

### Step 3: Export to ONNX

```bash
python scripts/export_onnx.py
```

Exports to ONNX format + quantizes for smaller size. Verifies the model works.

### Step 4: Deploy

```bash
cp models/volo-intent-quantized.onnx ../volo-api/models/intent.onnx
cp -r models/tokenizer ../volo-api/models/tokenizer
cp models/model-config.json ../volo-api/models/
```

## Project Structure

```
volo-model/
├── data/
│   ├── train.json          Generated training data
│   └── eval.json           Generated evaluation data
├── models/
│   ├── volo-intent/        HuggingFace checkpoint (after training)
│   ├── volo-intent.onnx    Full ONNX export
│   ├── volo-intent-quantized.onnx  Quantized (deploy this)
│   ├── tokenizer/          Tokenizer files for Go
│   └── model-config.json   Labels + config
├── scripts/
│   ├── prepare_dataset.py  Dataset preparation
│   ├── train.py            Model training
│   ├── export_onnx.py      ONNX export + quantization
│   └── evaluate.py         Evaluation + single-input testing
├── notebooks/              Jupyter notebooks for exploration
├── Dockerfile              Python environment
├── requirements.txt        Dependencies
└── README.md               This file
```

## Model Details

| Property | Value |
|----------|-------|
| Base model | distilbert-base-uncased (67M params) |
| Task | Sequence classification (5 intents) |
| Max input length | 64 tokens |
| Training time | ~10 min (CPU) / ~2 min (GPU) |
| ONNX size (quantized) | ~65 MB |
| Inference time | ~8ms (CPU) |

## Intent Labels

| Label | Description | Example |
|-------|-------------|---------|
| search | Web search | "search react hooks" |
| navigate | Open a site | "open youtube" |
| open-and-search | Open site + search | "open youtube lofi" |
| open-and-play | Open site + play | "open youtube play lofi" |
| browser-control | Tab/nav control | "close tab", "go back" |

## Evaluate a Single Input

```bash
python scripts/evaluate.py --input "play some jazz on spotify"
```

Output:
```json
{
  "text": "play some jazz on spotify",
  "intent": "open-and-play",
  "confidence": 0.91,
  "all_scores": { ... }
}
```
