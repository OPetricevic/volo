"""
Evaluate the trained model on custom inputs or the eval set.

Usage:
    # Evaluate on eval set
    python scripts/evaluate.py

    # Test single input
    python scripts/evaluate.py --input "open youtube play lofi"

    # Test multiple inputs from file
    python scripts/evaluate.py --file inputs.txt
"""

import argparse
import json
import sys
from pathlib import Path

import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer

MODEL_DIR = Path("/workspace/models") if Path("/workspace").exists() else Path("models")
DATA_DIR = Path("/workspace/data") if Path("/workspace").exists() else Path("data")

LABELS = ["search", "navigate", "open-and-search", "open-and-play", "browser-control"]


def load_model():
    """Load ONNX model and tokenizer."""
    model_path = MODEL_DIR / "volo-intent-quantized.onnx"
    if not model_path.exists():
        model_path = MODEL_DIR / "volo-intent.onnx"
    if not model_path.exists():
        print("Error: No ONNX model found. Run train.py and export_onnx.py first.")
        sys.exit(1)

    session = ort.InferenceSession(str(model_path))

    tokenizer_path = MODEL_DIR / "volo-intent"
    if not tokenizer_path.exists():
        tokenizer_path = MODEL_DIR / "tokenizer"
    tokenizer = AutoTokenizer.from_pretrained(str(tokenizer_path))

    return session, tokenizer


def predict(session, tokenizer, text):
    """Run inference on a single text."""
    inputs = tokenizer(
        text,
        return_tensors="np",
        max_length=64,
        padding="max_length",
        truncation=True,
    )

    outputs = session.run(None, {
        "input_ids": inputs["input_ids"],
        "attention_mask": inputs["attention_mask"],
    })

    logits = outputs[0][0]
    probs = np.exp(logits) / np.sum(np.exp(logits))
    predicted_idx = int(np.argmax(probs))

    return {
        "text": text,
        "intent": LABELS[predicted_idx],
        "confidence": float(probs[predicted_idx]),
        "all_scores": {label: float(probs[i]) for i, label in enumerate(LABELS)},
    }


def evaluate_dataset(session, tokenizer):
    """Evaluate on the full eval set."""
    eval_path = DATA_DIR / "eval.json"
    if not eval_path.exists():
        print("Error: eval.json not found. Run prepare_dataset.py first.")
        sys.exit(1)

    with open(eval_path) as f:
        eval_data = json.load(f)

    correct = 0
    total = len(eval_data)
    errors = []

    for item in eval_data:
        result = predict(session, tokenizer, item["text"])
        if result["intent"] == item["intent"]:
            correct += 1
        else:
            errors.append({
                "text": item["text"],
                "expected": item["intent"],
                "predicted": result["intent"],
                "confidence": result["confidence"],
            })

    accuracy = correct / total
    print(f"\n=== Evaluation Results ===")
    print(f"Accuracy: {accuracy:.4f} ({correct}/{total})")
    print(f"Errors: {len(errors)}")

    if errors:
        print(f"\n--- Sample Errors (first 10) ---")
        for err in errors[:10]:
            print(f"  '{err['text']}'")
            print(f"    Expected: {err['expected']}, Got: {err['predicted']} ({err['confidence']:.2f})")


def main():
    parser = argparse.ArgumentParser(description="Evaluate Volo intent model")
    parser.add_argument("--input", type=str, help="Single text to classify")
    parser.add_argument("--file", type=str, help="File with one text per line")
    args = parser.parse_args()

    session, tokenizer = load_model()

    if args.input:
        result = predict(session, tokenizer, args.input)
        print(json.dumps(result, indent=2))
    elif args.file:
        with open(args.file) as f:
            for line in f:
                line = line.strip()
                if line:
                    result = predict(session, tokenizer, line)
                    print(f"{result['intent']:20s} ({result['confidence']:.2f})  {line}")
    else:
        evaluate_dataset(session, tokenizer)


if __name__ == "__main__":
    main()
