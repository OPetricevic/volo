# Dataset

## Sources

The training dataset is assembled from two sources:

### 1. SNIPS Built-in Intents (MIT License)

A well-known NLU benchmark dataset. We filter and remap relevant intents:

| SNIPS Intent | → Volo Intent | Example |
|-------------|---------------|---------|
| SearchCreativeWork | search | "find me the movie with tom hanks" |
| SearchScreeningEvent | search | "find movie times near me" |
| PlayMusic | open-and-play | "play some rock music" |
| GetWeather | search | "what's the weather in london" |
| BookRestaurant | search | "book a table for two" |
| AddToPlaylist | open-and-search | "add this song to my playlist" |

### 2. Custom Browser Commands (hand-written + generated)

Specific to Volo's use case — browser voice commands:

- Navigate commands: "open youtube", "go to github", "take me to gmail"
- Search commands: "search for react hooks", "look up weather"
- Open + search: "open youtube lofi beats", "find jazz on spotify"
- Open + play: "open youtube play lofi", "play jazz on spotify"
- Browser control: "go back", "close tab", "scroll down", "new tab"

### 3. Augmentation

~20% of examples get filler words added (simulating real speech recognition output):
- Prefixes: "um", "uh", "like", "so", "okay", "hey", "well"
- Suffixes: "please", "can you"

## Format

```json
[
  {
    "text": "open youtube play lofi hip hop",
    "intent": "open-and-play",
    "source": "custom"
  },
  {
    "text": "um search for react hooks",
    "intent": "search",
    "source": "augmented"
  }
]
```

## Class Balance

The dataset is balanced to max 400 examples per class:

| Intent | Target count |
|--------|-------------|
| search | ~400 |
| navigate | ~400 |
| open-and-search | ~400 |
| open-and-play | ~400 |
| browser-control | ~50 (fewer natural variations) |

## Train/Eval Split

- 85% train, 15% eval
- Stratified by intent (proportional representation in both splits)
- Random seed: 42 (reproducible)

## Regenerating

```bash
python scripts/prepare_dataset.py
```

This downloads SNIPS, generates custom examples, augments, balances, and saves to `data/train.json` + `data/eval.json`.
