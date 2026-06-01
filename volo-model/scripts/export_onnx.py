"""
Export the trained DistilBERT model to ONNX format.

Usage:
    python scripts/export_onnx.py

Outputs:
    models/volo-intent.onnx           — ONNX model file
    models/volo-intent-quantized.onnx — Quantized (smaller, faster)
    models/tokenizer.json             — Tokenizer config for Go
"""

from pathlib import Path
from optimum.onnxruntime import ORTModelForSequenceClassification
from optimum.onnxruntime.configuration import AutoQuantizationConfig
from optimum.onnxruntime import ORTQuantizer
from transformers import AutoTokenizer
import onnxruntime as ort
import json
import shutil

MODEL_DIR = Path("/workspace/models/volo-intent") if Path("/workspace").exists() else Path("models/volo-intent")
OUTPUT_DIR = Path("/workspace/models") if Path("/workspace").exists() else Path("models")

LABELS = ["search", "navigate", "open-and-search", "open-and-play", "browser-control"]


def export_onnx():
    """Export HuggingFace model to ONNX."""
    print("=== ONNX Export ===\n")

    onnx_dir = OUTPUT_DIR / "onnx-export"
    onnx_dir.mkdir(parents=True, exist_ok=True)

    print(f"Loading model from: {MODEL_DIR}")
    model = ORTModelForSequenceClassification.from_pretrained(
        MODEL_DIR,
        export=True,
    )
    model.save_pretrained(str(onnx_dir))

    # Copy the model file to a clean name
    onnx_file = onnx_dir / "model.onnx"
    output_file = OUTPUT_DIR / "volo-intent.onnx"
    if onnx_file.exists():
        shutil.copy(onnx_file, output_file)
        size_mb = output_file.stat().st_size / (1024 * 1024)
        print(f"\nExported: {output_file} ({size_mb:.1f} MB)")

    return onnx_dir


def quantize(onnx_dir):
    """Quantize the ONNX model for smaller size and faster inference."""
    print("\n=== Quantization ===\n")

    quantizer = ORTQuantizer.from_pretrained(onnx_dir)
    qconfig = AutoQuantizationConfig.avx512_vnni(is_static=False, per_channel=False)

    quantized_dir = OUTPUT_DIR / "onnx-quantized"
    quantized_dir.mkdir(parents=True, exist_ok=True)

    quantizer.quantize(
        save_dir=str(quantized_dir),
        quantization_config=qconfig,
    )

    # Copy quantized model
    quantized_file = quantized_dir / "model_quantized.onnx"
    output_file = OUTPUT_DIR / "volo-intent-quantized.onnx"
    if quantized_file.exists():
        shutil.copy(quantized_file, output_file)
        size_mb = output_file.stat().st_size / (1024 * 1024)
        print(f"Quantized: {output_file} ({size_mb:.1f} MB)")


def export_tokenizer_config():
    """Export tokenizer config in a format the Go backend can use."""
    print("\n=== Tokenizer Config ===\n")

    tokenizer = AutoTokenizer.from_pretrained(str(MODEL_DIR))

    # Save the full tokenizer (Go can use tokenizer.json from HuggingFace format)
    tokenizer.save_pretrained(str(OUTPUT_DIR / "tokenizer"))

    # Also save a simplified config for reference
    config = {
        "model_name": "distilbert-base-uncased",
        "max_length": 64,
        "labels": LABELS,
        "label2id": {label: i for i, label in enumerate(LABELS)},
        "id2label": {str(i): label for i, label in enumerate(LABELS)},
        "vocab_size": tokenizer.vocab_size,
    }

    config_path = OUTPUT_DIR / "model-config.json"
    with open(config_path, "w") as f:
        json.dump(config, f, indent=2)
    print(f"Config: {config_path}")


def verify_onnx():
    """Quick verification that the ONNX model loads and runs."""
    print("\n=== Verification ===\n")

    model_path = OUTPUT_DIR / "volo-intent-quantized.onnx"
    if not model_path.exists():
        model_path = OUTPUT_DIR / "volo-intent.onnx"

    session = ort.InferenceSession(str(model_path))

    # Test input
    tokenizer = AutoTokenizer.from_pretrained(str(MODEL_DIR))
    inputs = tokenizer("open youtube play lofi", return_tensors="np", max_length=64, padding="max_length", truncation=True)

    outputs = session.run(None, {
        "input_ids": inputs["input_ids"],
        "attention_mask": inputs["attention_mask"],
    })

    import numpy as np
    logits = outputs[0][0]
    probs = np.exp(logits) / np.sum(np.exp(logits))
    predicted_idx = np.argmax(probs)
    predicted_label = LABELS[predicted_idx]
    confidence = probs[predicted_idx]

    print(f"Input: 'open youtube play lofi'")
    print(f"Predicted: {predicted_label} (confidence: {confidence:.3f})")
    print(f"\nAll probabilities:")
    for i, label in enumerate(LABELS):
        print(f"  {label}: {probs[i]:.3f}")

    print(f"\n✓ ONNX model works correctly")


def main():
    onnx_dir = export_onnx()
    quantize(onnx_dir)
    export_tokenizer_config()
    verify_onnx()

    print("\n=== Done ===")
    print(f"\nFiles ready for Go backend:")
    print(f"  {OUTPUT_DIR / 'volo-intent-quantized.onnx'}")
    print(f"  {OUTPUT_DIR / 'tokenizer/'}")
    print(f"  {OUTPUT_DIR / 'model-config.json'}")
    print(f"\nCopy these to volo-api/models/ for deployment.")


if __name__ == "__main__":
    main()
