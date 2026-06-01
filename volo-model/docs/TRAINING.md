# Training

## Model Choice: DistilBERT

| Property | Value |
|----------|-------|
| Architecture | DistilBERT (6 layers, 12 heads, 768 hidden) |
| Parameters | 67M |
| Pre-trained on | English Wikipedia + BookCorpus |
| Why this model | Small enough for CPU inference, large enough for good accuracy |

Alternatives considered:
- **BERT-base** (110M) — too large for fast CPU inference
- **TinyBERT** (14M) — too small, accuracy drops significantly
- **Phi-3** (3.8B) — overkill for classification, needs GPU

## Hyperparameters

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| Epochs | 8 (with early stopping) | Small dataset converges fast |
| Batch size | 32 | Fits in CPU memory comfortably |
| Learning rate | 3e-5 | Standard for BERT fine-tuning |
| Weight decay | 0.01 | Mild regularization |
| Warmup ratio | 0.1 | Prevents early divergence |
| Max length | 64 tokens | Voice commands are short (~10-20 words) |
| Early stopping | patience=3 | Stop if F1 doesn't improve for 3 epochs |

## Training Process

```
Epoch 1: lr warmup, model adapts to new task
Epoch 2-4: rapid improvement, F1 climbs
Epoch 5-8: fine-tuning, diminishing returns
Early stop: if no improvement for 3 consecutive epochs
```

Expected results:
- **Accuracy:** 92-96%
- **F1 Macro:** 90-94%
- **Training time:** ~10 min (CPU), ~2 min (GPU)

## Evaluation Metrics

| Metric | What it measures |
|--------|-----------------|
| Accuracy | % of correctly classified examples |
| F1 Macro | Average F1 across all classes (handles imbalance) |
| Per-class F1 | F1 for each intent separately |
| Confusion matrix | Which intents get confused with each other |

## ONNX Export

After training, the model is exported to ONNX for deployment in Go:

```
HuggingFace checkpoint (PyTorch)
    → ONNX export (via optimum library)
    → Dynamic quantization (INT8 weights)
    → Final model (~65MB)
```

### Quantization

We use dynamic quantization (INT8):
- Reduces model size by ~4x (250MB → 65MB)
- Inference speed improves ~2x
- Accuracy loss: <0.5% (negligible)

## Retraining

When to retrain:
- New intent categories added
- Accuracy drops on real-world commands (logged via audit)
- Significant new training data available (>100 new examples)

Process:
1. Add new examples to `data/` (or regenerate with updated script)
2. Run `python scripts/train.py`
3. Run `python scripts/export_onnx.py`
4. Copy ONNX to `volo-api/models/`
5. Restart API (`docker compose restart volo-api`)

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Low accuracy on "browser-control" | Add more control examples to custom data |
| Confuses "navigate" with "open-and-search" | Add more single-word site commands |
| Slow training | Reduce batch size, or use GPU container |
| OOM during training | Reduce max_length to 32, reduce batch size |
| ONNX export fails | Ensure optimum version matches transformers version |
