"""
Fine-tune DistilBERT for Volo intent classification.

Usage:
    python scripts/train.py

Outputs:
    models/volo-intent/       — HuggingFace model checkpoint
    models/volo-intent.onnx   — ONNX export for Go backend
"""

import json
import numpy as np
from pathlib import Path

import torch
from datasets import Dataset
from transformers import (
    AutoTokenizer,
    AutoModelForSequenceClassification,
    TrainingArguments,
    Trainer,
    EarlyStoppingCallback,
)
from sklearn.metrics import accuracy_score, f1_score, classification_report

# ─── Config ───────────────────────────────────────────────

MODEL_NAME = "distilbert-base-uncased"
OUTPUT_DIR = Path("/workspace/models/volo-intent") if Path("/workspace").exists() else Path("models/volo-intent")
DATA_DIR = Path("/workspace/data") if Path("/workspace").exists() else Path("data")

LABELS = ["search", "navigate", "open-and-search", "open-and-play", "browser-control"]
LABEL2ID = {label: i for i, label in enumerate(LABELS)}
ID2LABEL = {i: label for i, label in enumerate(LABELS)}

TRAINING_ARGS = {
    "num_train_epochs": 8,
    "per_device_train_batch_size": 32,
    "per_device_eval_batch_size": 64,
    "learning_rate": 3e-5,
    "weight_decay": 0.01,
    "warmup_ratio": 0.1,
    "eval_strategy": "epoch",
    "save_strategy": "epoch",
    "load_best_model_at_end": True,
    "metric_for_best_model": "f1_macro",
    "greater_is_better": True,
    "logging_steps": 50,
    "fp16": torch.cuda.is_available(),
}

MAX_LENGTH = 64  # Max tokens (voice commands are short)


# ─── Data Loading ─────────────────────────────────────────

def load_data():
    """Load train and eval JSON files."""
    with open(DATA_DIR / "train.json") as f:
        train_data = json.load(f)
    with open(DATA_DIR / "eval.json") as f:
        eval_data = json.load(f)

    print(f"Loaded {len(train_data)} train, {len(eval_data)} eval examples")
    return train_data, eval_data


def prepare_dataset(data, tokenizer):
    """Convert raw data to HuggingFace Dataset with tokenized inputs."""
    texts = [item["text"] for item in data]
    labels = [LABEL2ID[item["intent"]] for item in data]

    encodings = tokenizer(
        texts,
        truncation=True,
        padding="max_length",
        max_length=MAX_LENGTH,
        return_tensors="pt",
    )

    dataset = Dataset.from_dict({
        "input_ids": encodings["input_ids"],
        "attention_mask": encodings["attention_mask"],
        "labels": labels,
    })
    dataset.set_format("torch")
    return dataset


# ─── Metrics ──────────────────────────────────────────────

def compute_metrics(eval_pred):
    """Compute accuracy and macro F1."""
    logits, labels = eval_pred
    predictions = np.argmax(logits, axis=-1)
    return {
        "accuracy": accuracy_score(labels, predictions),
        "f1_macro": f1_score(labels, predictions, average="macro"),
    }


# ─── Main ─────────────────────────────────────────────────

def main():
    print("=== Volo Intent Classifier Training ===\n")
    print(f"Model: {MODEL_NAME}")
    print(f"Labels: {LABELS}")
    print(f"Output: {OUTPUT_DIR}\n")

    # Load data
    train_data, eval_data = load_data()

    # Tokenizer
    print("Loading tokenizer...")
    tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)

    # Prepare datasets
    print("Tokenizing...")
    train_dataset = prepare_dataset(train_data, tokenizer)
    eval_dataset = prepare_dataset(eval_data, tokenizer)

    # Model
    print("Loading model...")
    model = AutoModelForSequenceClassification.from_pretrained(
        MODEL_NAME,
        num_labels=len(LABELS),
        id2label=ID2LABEL,
        label2id=LABEL2ID,
    )

    # Training
    print("Starting training...\n")
    training_args = TrainingArguments(
        output_dir=str(OUTPUT_DIR),
        **TRAINING_ARGS,
    )

    trainer = Trainer(
        model=model,
        args=training_args,
        train_dataset=train_dataset,
        eval_dataset=eval_dataset,
        compute_metrics=compute_metrics,
        callbacks=[EarlyStoppingCallback(early_stopping_patience=3)],
    )

    trainer.train()

    # Save best model
    trainer.save_model(str(OUTPUT_DIR))
    tokenizer.save_pretrained(str(OUTPUT_DIR))

    # Final evaluation
    print("\n=== Final Evaluation ===")
    results = trainer.evaluate()
    print(f"Accuracy: {results['eval_accuracy']:.4f}")
    print(f"F1 Macro: {results['eval_f1_macro']:.4f}")

    # Detailed classification report
    predictions = trainer.predict(eval_dataset)
    preds = np.argmax(predictions.predictions, axis=-1)
    labels = predictions.label_ids
    print("\n" + classification_report(labels, preds, target_names=LABELS))

    print(f"\nModel saved to: {OUTPUT_DIR}")
    print("Next step: python scripts/export_onnx.py")


if __name__ == "__main__":
    main()
