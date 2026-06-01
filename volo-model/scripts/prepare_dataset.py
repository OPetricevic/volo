"""
Dataset preparation for Volo intent classification.

Combines:
1. SNIPS dataset (filtered to relevant intents, remapped to Volo labels)
2. Custom browser command examples (hand-written)
3. Generated variations (augmentation)

Output: data/train.json, data/eval.json
"""

import json
import random
from pathlib import Path
from datasets import load_dataset

# Volo intent labels
LABELS = ["search", "navigate", "open-and-search", "open-and-play", "browser-control"]

# SNIPS intent → Volo intent mapping
SNIPS_MAPPING = {
    "SearchCreativeWork": "search",
    "SearchScreeningEvent": "search",
    "PlayMusic": "open-and-play",
    "GetWeather": "search",
    "BookRestaurant": "search",
    "AddToPlaylist": "open-and-search",
    "RateBook": "search",
}

# Sites for variation generation
SITES = ["youtube", "spotify", "github", "reddit", "twitter", "gmail", "netflix", "twitch"]
QUERIES = [
    "lofi beats", "react hooks tutorial", "python machine learning",
    "jazz music", "funny cats", "cooking recipes", "workout playlist",
    "news today", "golang concurrency", "typescript generics",
    "best restaurants", "travel vlog", "coding interview prep",
    "ambient music", "game reviews", "tech news", "math tutorial",
]


def load_snips_data():
    """Load SNIPS dataset and remap to Volo intents."""
    print("Loading SNIPS dataset...")
    dataset = load_dataset("snips_built_in_intents", trust_remote_code=True)

    examples = []
    for split in ["train", "test"]:
        if split in dataset:
            for item in dataset[split]:
                snips_intent = item["label"]
                # Map label index to name
                label_name = dataset[split].features["label"].names[snips_intent]
                if label_name in SNIPS_MAPPING:
                    volo_intent = SNIPS_MAPPING[label_name]
                    examples.append({
                        "text": item["text"].lower().strip(),
                        "intent": volo_intent,
                        "source": "snips",
                    })

    print(f"  Loaded {len(examples)} examples from SNIPS")
    return examples


def generate_custom_examples():
    """Hand-written examples specific to browser voice commands."""
    examples = []

    # ─── Search ───
    search_templates = [
        "search {query}",
        "search for {query}",
        "look up {query}",
        "find {query}",
        "google {query}",
        "{query}",  # bare query (fallback)
        "what is {query}",
        "how to {query}",
        "where can i find {query}",
    ]
    for template in search_templates:
        for query in QUERIES:
            examples.append({"text": template.format(query=query), "intent": "search", "source": "custom"})

    # ─── Navigate ───
    navigate_templates = [
        "open {site}",
        "go to {site}",
        "goto {site}",
        "take me to {site}",
        "launch {site}",
        "navigate to {site}",
    ]
    for template in navigate_templates:
        for site in SITES:
            examples.append({"text": template.format(site=site), "intent": "navigate", "source": "custom"})

    # ─── Open and Search ───
    open_search_templates = [
        "open {site} {query}",
        "search {site} for {query}",
        "find {query} on {site}",
        "look up {query} on {site}",
        "{site} {query}",
    ]
    for template in open_search_templates:
        for site in ["youtube", "spotify", "github", "reddit"]:
            for query in random.sample(QUERIES, min(5, len(QUERIES))):
                examples.append({"text": template.format(site=site, query=query), "intent": "open-and-search", "source": "custom"})

    # ─── Open and Play ───
    play_templates = [
        "open {site} play {query}",
        "play {query} on {site}",
        "put on {query} on {site}",
        "play some {query}",
        "play {query}",
    ]
    for template in play_templates:
        for site in ["youtube", "spotify"]:
            for query in random.sample(QUERIES, min(5, len(QUERIES))):
                examples.append({"text": template.format(site=site, query=query), "intent": "open-and-play", "source": "custom"})

    # ─── Browser Control ───
    control_examples = [
        ("go back", "browser-control"),
        ("back", "browser-control"),
        ("go forward", "browser-control"),
        ("forward", "browser-control"),
        ("close tab", "browser-control"),
        ("close this tab", "browser-control"),
        ("close the tab", "browser-control"),
        ("new tab", "browser-control"),
        ("open new tab", "browser-control"),
        ("open a new tab", "browser-control"),
        ("scroll down", "browser-control"),
        ("scroll up", "browser-control"),
        ("page down", "browser-control"),
        ("page up", "browser-control"),
        ("refresh", "browser-control"),
        ("reload", "browser-control"),
        ("reload the page", "browser-control"),
        ("stop", "browser-control"),
        ("stop loading", "browser-control"),
    ]
    for text, intent in control_examples:
        examples.append({"text": text, "intent": intent, "source": "custom"})

    print(f"  Generated {len(examples)} custom examples")
    return examples


def augment_with_filler(examples, ratio=0.2):
    """Add speech recognition artifacts: filler words, hesitations."""
    fillers = ["um", "uh", "like", "so", "okay", "hey", "well", "please", "can you"]
    augmented = []

    sample_size = int(len(examples) * ratio)
    for ex in random.sample(examples, min(sample_size, len(examples))):
        filler = random.choice(fillers)
        position = random.choice(["prefix", "suffix"])
        if position == "prefix":
            new_text = f"{filler} {ex['text']}"
        else:
            new_text = f"{ex['text']} {filler}"

        augmented.append({
            "text": new_text,
            "intent": ex["intent"],
            "source": "augmented",
        })

    print(f"  Augmented {len(augmented)} examples with filler words")
    return augmented


def balance_dataset(examples, max_per_class=300):
    """Balance classes so no single intent dominates."""
    by_intent = {}
    for ex in examples:
        by_intent.setdefault(ex["intent"], []).append(ex)

    balanced = []
    for intent, items in by_intent.items():
        if len(items) > max_per_class:
            balanced.extend(random.sample(items, max_per_class))
        else:
            balanced.extend(items)

    print(f"  Balanced to {len(balanced)} examples")
    for intent in LABELS:
        count = sum(1 for ex in balanced if ex["intent"] == intent)
        print(f"    {intent}: {count}")

    return balanced


def split_and_save(examples, train_ratio=0.85):
    """Split into train/eval and save as JSON."""
    random.shuffle(examples)
    split_idx = int(len(examples) * train_ratio)
    train = examples[:split_idx]
    eval_data = examples[split_idx:]

    output_dir = Path("/workspace/data") if Path("/workspace").exists() else Path("data")
    output_dir.mkdir(exist_ok=True)

    with open(output_dir / "train.json", "w") as f:
        json.dump(train, f, indent=2)

    with open(output_dir / "eval.json", "w") as f:
        json.dump(eval_data, f, indent=2)

    print(f"\nSaved: {len(train)} train, {len(eval_data)} eval")
    print(f"  → {output_dir / 'train.json'}")
    print(f"  → {output_dir / 'eval.json'}")


def main():
    random.seed(42)
    print("=== Volo Dataset Preparation ===\n")

    # 1. Load SNIPS
    snips = load_snips_data()

    # 2. Custom browser examples
    custom = generate_custom_examples()

    # 3. Combine
    all_examples = snips + custom

    # 4. Augment with filler words
    augmented = augment_with_filler(all_examples)
    all_examples.extend(augmented)

    # 5. Balance
    balanced = balance_dataset(all_examples, max_per_class=400)

    # 6. Split and save
    split_and_save(balanced)

    print("\nDone! Ready for training.")


if __name__ == "__main__":
    main()
