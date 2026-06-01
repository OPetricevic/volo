# Volo — Voice-First Browser Assistant

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
| Intent parsing | Rule-based + fuzzy matching (evolve to ML later) |
| Pattern model | Weighted frequency counters + recency decay |
| API style | REST (JSON) |
| Deployment | Docker / single binary |

### Volo Landing
| Layer | Choice |
|-------|--------|
| Framework | React + TypeScript |
| Bundler | Vite |
| Styling | Tailwind CSS |
| Hosting | Vercel / Netlify / static |
| Content | Download links, features, demo video placeholder |

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
│   │   ├── intent\          — Intent parser
│   │   ├── model\           — Pattern learning model
│   │   ├── storage\         — Database layer
│   │   └── middleware\      — Auth, logging, CORS
│   ├── migrations\          — SQL migrations
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
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
└── SPEC.md                  — This file
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

### Phase 3 — Pattern Learning
12. Implement frequency model in Go
13. Suggestions endpoint based on history
14. Time-of-day weighting
15. Fuzzy matching for site names

### Phase 4 — Desktop App
16. Scaffold Electron app with React renderer
17. System tray + global hotkey
18. Embedded BrowserView for web navigation
19. Share voice logic with extension
20. Desktop app launching (shell exec)

### Phase 5 — Landing Page
21. Scaffold React + Vite site
22. Hero section with demo/video placeholder
23. Features breakdown
24. Download links (extension store + desktop installer)
25. Deploy to static hosting

### Phase 6 — Polish
26. Visual feedback (listening animation, command confirmation)
27. Error handling and retry logic
28. Settings sync across devices
29. Onboarding flow (mic permission, tutorial)
30. Cross-browser testing (Chrome, Opera, Edge)

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
- Cloud deployment (local dev first)
- Real ML model training (start with frequency heuristics)
- Voice synthesis (Volo doesn't talk back — yet)
- Multi-user / team features
