# Volo — Voice-First Browser Assistant

> **Documentation:** [DOCS.md](./DOCS.md) | [DOCS.html](./DOCS.html) (visual)
> **Decisions:** [DECISIONS.md](./DECISIONS.md) (why things are the way they are)
> **Visual Spec:** [SPEC.html](./SPEC.html)
> **Desktop Spec:** [volo-desktop/SPEC.md](./volo-desktop/SPEC.md)

## Overview

Volo is a voice-activated personal assistant that lives in your browser and on your desktop. You say **"Hey Volo"** and it listens, then executes your command — searching the web, opening apps/sites, navigating within sites, or playing content. It learns your patterns over time using a lightweight AI model to surface your most likely intent faster.

---

## Products

| Product | Description | Tech |
|---------|-------------|------|
| **Volo Extension** | Chromium browser extension (Chrome, Opera, Edge) | React + Manifest V3 + Web Speech API |
| **Volo Desktop** | Standalone desktop app with embedded browser | Electron + React |
| **Volo API** | Backend service powering intent parsing + learning | Go |
| **Volo Landing** | Marketing/download page | React + TypeScript (Vite) |

All clients share the same core voice logic and communicate with the same Go backend.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      User Voice                          │
└──────────────────────────┬──────────────────────────────┘
                           │ "Hey Volo, ..."
                           ▼
┌──────────────────────────────────────────────────────────┐
│              Voice Recognition Layer                       │
│         (Web Speech API / always-listening)                │
│                                                           │
│  ┌─────────────┐    ┌──────────────────┐                 │
│  │  Extension   │    │  Desktop (Electron)│                │
│  │  (content    │    │  (renderer process)│                │
│  │   script)    │    │                    │                │
│  └──────┬───────┘    └────────┬───────────┘               │
└─────────┼─────────────────────┼──────────────────────────┘
          │                     │
          ▼                     ▼
┌──────────────────────────────────────────────────────────┐
│                    Volo API (Go)                           │
│                                                           │
│  ┌────────────┐  ┌────────────┐  ┌───────────────────┐   │
│  │   Intent    │  │   Search   │  │  Pattern Learning  │  │
│  │   Parser    │  │   History  │  │  (Frequency Model) │  │
│  └────────────┘  └────────────┘  └───────────────────┘   │
│                                                           │
│  ┌────────────────────────────────────────────────────┐   │
│  │              SQLite / PostgreSQL                     │   │
│  └────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────┘
```

---

## Core Features

### 1. Wake Word Detection

- **Trigger:** "Hey Volo" (always-listening mode)
- **Microphone permissions:** Asked on install with three options:
  - **Always** — continuously listens for wake word
  - **Once** — listens only when manually activated (click/hotkey)
  - **Off** — disabled entirely
- User can change this anytime in settings
- Visual indicator when Volo is listening (pulsing icon, overlay)

### 2. Voice Commands

| Command Pattern | Action |
|----------------|--------|
| "Hey Volo, search [query]" | Google search for the query |
| "Hey Volo, [query]" | Smart search (defaults to Google) |
| "Hey Volo, open [site]" | Navigate to site (youtube.com, github.com, etc.) |
| "Hey Volo, open YouTube [query]" | Open YouTube and search for query |
| "Hey Volo, open YouTube play [query]" | Open YouTube, search, and play first result |
| "Hey Volo, open Spotify [query]" | Open Spotify web and search |
| "Hey Volo, open [app]" | Open desktop app (desktop version only) |
| "Hey Volo, go back" | Browser back |
| "Hey Volo, scroll down/up" | Scroll the page |
| "Hey Volo, close tab" | Close current tab |
| "Hey Volo, new tab" | Open new tab |

### 3. Intent Parsing (Backend)

The Go backend receives the raw transcript and determines:
- **Action type:** search, navigate, open-and-search, open-and-play, browser-control
- **Target:** which site/app
- **Query:** the search terms
- **Confidence score:** how sure the model is

Example:
```
Input: "open youtube play lofi hip hop"
→ { action: "open-and-play", target: "youtube", query: "lofi hip hop", confidence: 0.95 }
```

### 4. Pattern Learning (AI Layer)

A lightweight model that improves over time:
- **Frequency tracking:** Most searched terms, most visited sites
- **Time-of-day patterns:** "In the morning, user usually opens Gmail first"
- **Autocomplete suggestions:** After "Hey Volo", suggest likely commands based on history
- **Fuzzy matching:** If user says "open tube" → likely means YouTube
- **Personalized ranking:** If user searches "react hooks" 10x, rank it higher in suggestions

Implementation: Weighted frequency model stored per-user in the database. Not a neural network — more like a smart trie + recency-weighted counters. Can evolve to something more sophisticated later.

### 5. Desktop App (Electron)

- Lives in system tray (like Discord)
- Always running in background
- Has its own embedded Chromium browser (BrowserView/WebContentsView)
- Can open URLs in its own window or in the user's default browser
- Global hotkey to activate (e.g., Ctrl+Shift+V)
- Shows a small overlay when processing voice
- Can launch desktop apps (via shell commands)

### 6. Browser Extension

- Manifest V3 (service worker + content scripts)
- Content script injects the voice listener on every page
- Popup shows status, recent commands, settings
- Badge icon changes color when listening (green = active, gray = off)
- Works on Chrome, Opera, Edge (all Chromium-based)

---

## Microphone Permission Flow

```
┌─────────────────────────────────────┐
│         First Install                │
│                                      │
│  "Volo needs microphone access       │
│   to listen for voice commands."     │
│                                      │
│  ┌──────────┐ ┌──────┐ ┌─────┐     │
│  │  Always   │ │ Once │ │ Off │     │
│  └──────────┘ └──────┘ └─────┘     │
│                                      │
│  "You can change this anytime        │
│   in Volo settings."                 │
└─────────────────────────────────────┘
```

- **Always:** Service worker keeps speech recognition active. Restarts on silence timeout.
- **Once:** User must click the Volo icon or press hotkey to activate listening.
- **Off:** Extension is passive. No mic access. Can still be used via text input.

---

## Tech Stack Detail

### Volo Extension
| Layer | Choice |
|-------|--------|
| UI | React 18+ with TypeScript |
| Bundler | Vite (with CRXJS or manual manifest) |
| Voice | Web Speech API (SpeechRecognition) |
| State | Zustand (lightweight) |
| Styling | Tailwind CSS |
| Manifest | V3 (service worker, content scripts) |
| Storage | chrome.storage.local for settings |
| Comms | chrome.runtime messaging between popup/content/background |

### Volo Desktop
| Layer | Choice |
|-------|--------|
| Framework | Electron |
| UI | React 18+ with TypeScript (shared with extension) |
| Voice | Web Speech API (in renderer) or native speech lib |
| Browser | Electron BrowserView / WebContentsView |
| Tray | Electron Tray API |
| Hotkey | Electron globalShortcut |
| IPC | Electron IPC (main ↔ renderer) |
| Auto-update | electron-updater |

### Volo API (Go)
| Layer | Choice |
|-------|--------|
| Language | Go 1.22+ |
| Framework | Chi or standard net/http |
| Database | SQLite (dev) / PostgreSQL (prod) |
| Auth | JWT or API key per device |
| Intent parsing | Rule-based (fast path) + DistilBERT ONNX (fallback) |
| Pattern model | Weighted frequency counters + recency decay |
| ML Runtime | ONNX Runtime (via onnxruntime-go) |
| API style | REST (JSON) |
| Deployment | Docker / single binary → VM |

### Volo Model (Python — training only)
| Layer | Choice |
|-------|--------|
| Base model | DistilBERT (67M params) |
| Framework | HuggingFace Transformers + PyTorch |
| Task | Sequence classification (intent) + token classification (slots) |
| Export | ONNX format (portable, runs in Go) |
| Dataset | Custom ~500-1000 labeled voice commands |
| Training | CPU-friendly, ~10 min on any machine |

### Volo Landing
| Layer | Choice |
|-------|--------|
| Framework | React + TypeScript |
| Bundler | Vite |
| Styling | Tailwind CSS |
| Hosting | Vercel / Netlify / static |
| Content | Download links, features, demo video placeholder |

---

## AI / Intent Classification Pipeline

### Architecture: Hybrid (Rule-Based + Fine-Tuned Model)

```
Transcript: "open youtube play lofi hip hop"
                    │
                    ▼
    ┌───────────────────────────────┐
    │     Rule-Based Parser          │  ← Fast path (~1ms)
    │     Pattern matching on        │
    │     keywords: open, search,    │
    │     play, go, close, new...    │
    │                                │
    │     confidence = 0.92 ✓        │
    └───────────────┬───────────────┘
                    │
            if confidence < 0.80
                    │
                    ▼
    ┌───────────────────────────────┐
    │     DistilBERT (ONNX)          │  ← Fallback (~8ms)
    │     Fine-tuned classifier      │
    │     Handles ambiguous/novel    │
    │     commands                   │
    └───────────────────────────────┘
```

### Why DistilBERT?

- **67M parameters** — tiny, runs on CPU in milliseconds
- **No GPU required** — deploys on any VM or server
- **ONNX export** — portable binary, loads directly in Go
- **Fine-tunable** — train on your own dataset in minutes
- **Proven** — well-documented, battle-tested architecture

### Model Tasks

The model handles two tasks simultaneously:

1. **Intent Classification** (what action to take)
   - Labels: `search`, `navigate`, `open-and-search`, `open-and-play`, `browser-control`

2. **Slot Extraction** (what entities are in the command)
   - Slots: `target` (youtube, github, etc.), `query` (search terms)

### Training Pipeline

```
┌─────────────────────────────────────────────────────────┐
│  1. Create Dataset                                       │
│     - 500-1000 labeled commands (JSON/CSV)               │
│     - Mix of simple + ambiguous examples                 │
│     - Include user variations ("open tube" = youtube)    │
└──────────────────────────┬──────────────────────────────┘
                           ▼
┌─────────────────────────────────────────────────────────┐
│  2. Fine-Tune (Python + HuggingFace)                     │
│     - Load pretrained distilbert-base-uncased            │
│     - Add classification head for intents                │
│     - Train ~5 epochs, batch size 16                     │
│     - Takes ~10 min on CPU                               │
└──────────────────────────┬──────────────────────────────┘
                           ▼
┌─────────────────────────────────────────────────────────┐
│  3. Export to ONNX                                       │
│     - torch.onnx.export() → model.onnx                  │
│     - Tokenizer config → tokenizer.json                  │
│     - Single portable file (~250MB → quantized ~65MB)    │
└──────────────────────────┬──────────────────────────────┘
                           ▼
┌─────────────────────────────────────────────────────────┐
│  4. Load in Go Backend                                   │
│     - onnxruntime-go loads model.onnx at startup         │
│     - Tokenize input → run inference → get labels        │
│     - ~8ms per inference on CPU                          │
└─────────────────────────────────────────────────────────┘
```

### Example Training Data (dataset format)

```json
[
  { "text": "search for react hooks", "intent": "search", "target": null, "query": "react hooks" },
  { "text": "open youtube", "intent": "navigate", "target": "youtube.com", "query": null },
  { "text": "open youtube play lofi", "intent": "open-and-play", "target": "youtube.com", "query": "lofi" },
  { "text": "find that golang tutorial", "intent": "search", "target": null, "query": "golang tutorial" },
  { "text": "go to github", "intent": "navigate", "target": "github.com", "query": null },
  { "text": "close this tab", "intent": "browser-control", "target": null, "query": null },
  { "text": "play some jazz on spotify", "intent": "open-and-play", "target": "spotify.com", "query": "jazz" }
]
```

### Go Integration

```go
package intent

import (
    ort "github.com/yalue/onnxruntime_go"
)

type ModelParser struct {
    session   *ort.Session
    tokenizer *Tokenizer  // custom tokenizer matching HuggingFace output
}

func (m *ModelParser) Parse(transcript string) (CommandResult, error) {
    // 1. Tokenize
    tokens := m.tokenizer.Encode(transcript)
    
    // 2. Run inference
    output, err := m.session.Run(tokens)
    if err != nil {
        return CommandResult{}, err
    }
    
    // 3. Decode labels
    intent := decodeIntent(output[0])   // "open-and-play"
    slots := decodeSlots(output[1])     // {target: "youtube.com", query: "lofi"}
    
    return CommandResult{
        Action:     intent,
        Target:     slots.Target,
        Query:      slots.Query,
        Confidence: output[0].MaxScore(),
    }, nil
}
```

### Continuous Improvement

The model improves over time without retraining:
1. **Frequency model** handles personalization (most-used commands rank higher)
2. **Correction logging** — if user repeats a command differently, log it as training data
3. **Periodic retraining** — batch new examples, retrain model, hot-swap ONNX file

---

## Deployment

### Target: Self-Hosted VM

The entire backend deploys as a single Docker container on your VM.

```
┌─────────────────────────────────────────────┐
│              Your VM                         │
│                                              │
│  ┌────────────────────────────────────────┐  │
│  │         Docker Container                │  │
│  │                                         │  │
│  │  ┌─────────────┐  ┌────────────────┐   │  │
│  │  │  Go Binary   │  │  model.onnx    │   │  │
│  │  │  (volo-api)  │  │  (65MB)        │   │  │
│  │  └──────┬───────┘  └────────────────┘   │  │
│  │         │                                │  │
│  │  ┌──────┴───────┐                       │  │
│  │  │  SQLite DB    │  (or Postgres)        │  │
│  │  │  (data.db)    │                       │  │
│  │  └──────────────┘                       │  │
│  └────────────────────────────────────────┘  │
│                                              │
│  ┌────────────────────────────────────────┐  │
│  │  Caddy / Nginx (reverse proxy + TLS)    │  │
│  │  api.volo.yourdomain.com → :8080        │  │
│  └────────────────────────────────────────┘  │
│                                              │
│  ┌────────────────────────────────────────┐  │
│  │  Volo Landing (static files)            │  │
│  │  volo.yourdomain.com → /var/www/volo    │  │
│  └────────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

### Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o volo-api ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates libc6-compat
WORKDIR /app
COPY --from=builder /app/volo-api .
COPY --from=builder /app/models/ ./models/
COPY --from=builder /app/migrations/ ./migrations/
EXPOSE 8080
CMD ["./volo-api"]
```

### docker-compose.yml

```yaml
version: "3.8"
services:
  volo-api:
    build: ./volo-api
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data          # SQLite persistence
      - ./models:/app/models      # ONNX model files
    environment:
      - VOLO_DB_PATH=/app/data/volo.db
      - VOLO_MODEL_PATH=/app/models/intent.onnx
      - VOLO_JWT_SECRET=${JWT_SECRET}
      - VOLO_PORT=8080
    restart: unless-stopped

  caddy:
    image: caddy:2-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - ./volo-landing/dist:/var/www/landing
      - caddy_data:/data
    restart: unless-stopped

volumes:
  caddy_data:
```

### Caddyfile

```
api.volo.yourdomain.com {
    reverse_proxy volo-api:8080
}

volo.yourdomain.com {
    root * /var/www/landing
    file_server
}
```

### Deploy Commands

```bash
# On your VM:
git clone <your-repo>
cd volo-api

# Build and start
docker compose up -d --build

# Update model (hot-swap without restart)
scp models/intent.onnx vm:/path/to/models/
docker compose restart volo-api

# View logs
docker compose logs -f volo-api
```

### Requirements

- **VM specs:** 1 vCPU, 1GB RAM minimum (DistilBERT ONNX is lightweight)
- **Storage:** ~500MB (Go binary + model + DB)
- **OS:** Any Linux (Ubuntu/Debian recommended)
- **Ports:** 80, 443 (Caddy handles TLS automatically via Let's Encrypt)

---

## API Endpoints (Go Backend)

```
POST   /api/v1/command          — Process a voice command
GET    /api/v1/suggestions      — Get autocomplete suggestions based on history
GET    /api/v1/history          — Get recent command history
DELETE /api/v1/history          — Clear history

POST   /api/v1/auth/register    — Register device/user
POST   /api/v1/auth/token       — Get auth token

GET    /api/v1/settings         — Get user settings (synced across devices)
PUT    /api/v1/settings         — Update settings

GET    /api/v1/health           — Health check
```

### Command Request/Response

```json
// POST /api/v1/command
// Request:
{
  "transcript": "open youtube play lofi hip hop",
  "context": {
    "current_url": "https://google.com",
    "timestamp": "2026-06-01T08:35:00Z"
  }
}

// Response:
{
  "action": "open-and-play",
  "target": "youtube.com",
  "query": "lofi hip hop",
  "confidence": 0.95,
  "suggestions": ["lofi hip hop radio", "lofi hip hop beats"],
  "execute": {
    "url": "https://www.youtube.com/results?search_query=lofi+hip+hop",
    "auto_play": true
  }
}
```

---

## Database Schema

```sql
-- Users / devices
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ DEFAULT now(),
  device_name TEXT,
  settings JSONB DEFAULT '{}'
);

-- Command history
CREATE TABLE commands (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  transcript TEXT NOT NULL,
  parsed_action TEXT NOT NULL,
  parsed_target TEXT,
  parsed_query TEXT,
  confidence REAL,
  executed_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_commands_user_time ON commands(user_id, executed_at DESC);

-- Pattern learning (frequency model)
CREATE TABLE patterns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  pattern_type TEXT NOT NULL,  -- 'search', 'site', 'command'
  value TEXT NOT NULL,
  frequency INTEGER DEFAULT 1,
  last_used TIMESTAMPTZ DEFAULT now(),
  hour_weights JSONB DEFAULT '{}',  -- { "8": 5, "9": 3 } = used 5x at 8am
  UNIQUE(user_id, pattern_type, value)
);

CREATE INDEX idx_patterns_user_freq ON patterns(user_id, frequency DESC);
```

---

## Project Structure

```
c:\Learn\Thesis\
├── volo-extension\          — Browser extension (React + Manifest V3)
│   ├── src\
│   │   ├── background\      — Service worker
│   │   ├── content\         — Content script (voice listener)
│   │   ├── popup\           — Popup UI (React)
│   │   ├── options\         — Settings page (React)
│   │   ├── shared\          — Shared types, API client, voice logic
│   │   └── assets\          — Icons, sounds
│   ├── public\
│   │   └── manifest.json
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── package.json
│
├── volo-desktop\            — Electron desktop app
│   ├── src\
│   │   ├── main\            — Electron main process
│   │   ├── renderer\        — React UI (shared with extension)
│   │   ├── preload\         — Preload scripts
│   │   └── shared\          — Shared voice/API logic
│   ├── electron-builder.yml
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── package.json
│
├── volo-api\                — Go backend
│   ├── cmd\
│   │   └── server\
│   │       └── main.go
│   ├── internal\
│   │   ├── handler\         — HTTP handlers
│   │   ├── intent\          — Intent parser (rule-based + ONNX)
│   │   ├── model\           — Pattern learning model
│   │   ├── storage\         — Database layer
│   │   └── middleware\      — Auth, logging, CORS
│   ├── models\              — ONNX model files (intent.onnx, tokenizer.json)
│   ├── migrations\          — SQL migrations
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── Makefile
│
├── volo-model\              — ML training pipeline (Python)
│   ├── data\
│   │   ├── train.json       — Training dataset
│   │   └── eval.json        — Evaluation dataset
│   ├── scripts\
│   │   ├── train.py         — Fine-tune DistilBERT
│   │   ├── export_onnx.py   — Export to ONNX format
│   │   ├── evaluate.py      — Test accuracy
│   │   └── generate_data.py — Helper to generate training examples
│   ├── models\              — Output: trained model + ONNX export
│   ├── requirements.txt     — torch, transformers, onnx, datasets
│   └── README.md
│
├── volo-landing\            — Marketing site
│   ├── src\
│   │   ├── components\
│   │   ├── pages\
│   │   └── assets\
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── package.json
│
├── SPEC.md                  — This file
└── SPEC.html                — Visual spec (open in browser)
```

---

## Implementation Order

### Phase 1 — Core Voice Loop (Extension)
1. Scaffold extension with Vite + React + Manifest V3
2. Implement Web Speech API listener in content script
3. Wake word detection ("Hey Volo")
4. Basic command parsing (client-side, no backend yet)
5. Execute simple commands: search Google, open URL
6. Popup UI: status indicator, mic toggle

### Phase 2 — Go Backend
7. Scaffold Go API with Chi router
8. Implement `/api/v1/command` endpoint with rule-based intent parser
9. Database setup (SQLite for dev)
10. Command history storage
11. Connect extension to backend API

### Phase 3 — ML Model + Pattern Learning
12. Create training dataset (~500-1000 labeled commands)
13. Fine-tune DistilBERT on intent classification
14. Export to ONNX, integrate into Go backend via onnxruntime-go
15. Implement frequency model for personalization
16. Suggestions endpoint based on history + time-of-day
17. Fuzzy matching for site names

### Phase 4 — Desktop App
18. Scaffold Electron app with React renderer
19. System tray + global hotkey
20. Embedded BrowserView for web navigation
21. Share voice logic with extension
22. Desktop app launching (shell exec)

### Phase 5 — Landing Page
23. Scaffold React + Vite site
24. Hero section with demo/video placeholder
25. Features breakdown
26. Download links (extension store + desktop installer)
27. Deploy to static hosting

### Phase 6 — Deployment + Polish
28. Dockerize Go backend + ONNX model
29. Deploy to VM (docker-compose + Caddy for TLS)
30. Visual feedback (listening animation, command confirmation)
31. Error handling and retry logic
32. Settings sync across devices
33. Onboarding flow (mic permission, tutorial)
34. Cross-browser testing (Chrome, Opera, Edge)

---

## Design Direction

- **Brand:** "Volo" — means "I fly" in Italian/Latin. Fast, lightweight, goes where you want.
- **Colors:** Dark UI with an accent color (electric blue or violet)
- **Popup/Overlay:** Minimal, non-intrusive. Small floating pill when listening.
- **Desktop:** Clean, modern. Dark sidebar + light content area.
- **Landing:** Bold hero, animated waveform visual, clear CTAs.

---

## Thesis Angles

This project supports several academic angles:

1. **Voice-driven HCI** — How voice interfaces change browser interaction patterns
2. **Personalized search** — Lightweight ML for user behavior prediction
3. **Cross-platform architecture** — Shared logic across extension + desktop
4. **Intent parsing** — NLP-lite approach to command understanding
5. **Edge computing** — Client-side wake word detection vs server-side processing

---

## Out of Scope (for now)

- Mobile app
- Firefox support (Manifest V2 differences)
- Multi-language voice recognition (English only initially)
- Voice synthesis (Volo doesn't talk back — yet)
- Multi-user / team features
- GPU inference (CPU-only is fine for DistilBERT)
- Real-time model retraining (batch retrain manually when needed)
