# Volo Desktop — Spec

## Overview

Electron desktop app with a Discord/Claude-like chat interface. Lives in system tray, activated via global hotkey. Users can speak commands (same as extension) OR chat with Volo about their search history using a local Gemma model via Ollama.

---

## Two Modes

### 1. Command Mode (same as extension)
- "Hey Volo, open YouTube play lofi" → executes immediately
- Opens embedded browser or default browser
- Same parser, same API, same behavior

### 2. Chat Mode (new — history queries)
- "What did I listen to yesterday?" → Gemma translates to DB query → formats response
- "Find that React thing from last week" → fuzzy search on history
- "What do I usually open in the morning?" → pattern analysis
- Conversational, but scoped to user's own data

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Electron App                                                │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Main Process                                            │ │
│  │  • Window management (tray, main window)                 │ │
│  │  • Global hotkey (Ctrl+Shift+V)                          │ │
│  │  • IPC bridge to renderer                                │ │
│  │  • Auto-updater                                          │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Renderer (React)                                        │ │
│  │  • Sidebar (nav, account, settings)                      │ │
│  │  • Chat view (messages, input, voice button)             │ │
│  │  • History view (searchable list)                        │ │
│  │  • Settings view (mic mode, account, theme)              │ │
│  │  • Embedded browser (BrowserView for "open" commands)    │ │
│  └─────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
         │                              │
         │ HTTP (same API)              │ HTTP (Ollama)
         ▼                              ▼
┌─────────────────┐          ┌─────────────────────┐
│  Volo API       │          │  Ollama (local)      │
│  (Go backend)   │          │  Gemma 2B            │
│  :8080          │          │  :11434              │
└─────────────────┘          └─────────────────────┘
```

---

## UI Layout

```
┌──────────────────────────────────────────────────────────────┐
│  ┌──────────┐  ┌──────────────────────────────────────────┐  │
│  │          │  │  Volo                            ─ □ ✕   │  │
│  │  [Logo]  │  ├──────────────────────────────────────────┤  │
│  │          │  │                                          │  │
│  │  Chat    │  │  ┌─────────────────────────────────────┐ │  │
│  │  History │  │  │ 🤖 Hey! Ask me about your history.  │ │  │
│  │  Browse  │  │  └─────────────────────────────────────┘ │  │
│  │          │  │                                          │  │
│  │          │  │  ┌─────────────────────────────────────┐ │  │
│  │          │  │  │ 👤 What did I listen to yesterday?  │ │  │
│  │          │  │  └─────────────────────────────────────┘ │  │
│  │          │  │                                          │  │
│  │          │  │  ┌─────────────────────────────────────┐ │  │
│  │          │  │  │ 🤖 Yesterday you played:            │ │  │
│  │          │  │  │    • lofi beats (YouTube, 9:14am)   │ │  │
│  │          │  │  │    • jazz playlist (Spotify, 2pm)   │ │  │
│  │          │  │  └─────────────────────────────────────┘ │  │
│  │          │  │                                          │  │
│  │ ──────── │  │  ┌──────────────────────────────┬─────┐ │  │
│  │ Settings │  │  │  Ask something...            │ 🎤  │ │  │
│  │ Account  │  │  └──────────────────────────────┴─────┘ │  │
│  └──────────┘  └──────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

---

## Color Palette (Discord-inspired, dark)

```
Background:     #1a1a2e (deep navy)
Surface:        #16213e (sidebar, cards)
Surface-2:      #0f3460 (hover, active states)
Border:         #2a2a4a
Text:           #e4e4f0
Text muted:     #8888a8
Accent:         #0ea5e9 (sky blue — matches brand)
Accent-2:       #14b8a6 (teal — secondary)
Success:        #34d399
Warning:        #fbbf24
Danger:         #f87171
```

---

## Features

### Sidebar
- Logo + app name
- Navigation: Chat, History, Browse, Settings
- Account section at bottom (avatar, name, logout)
- Collapse/expand on small screens

### Chat View
- Message bubbles (user = right/blue, volo = left/dark)
- Voice input button (big mic icon, pulses when listening)
- Text input (type or speak)
- Auto-scroll to latest message
- Loading indicator when Gemma is thinking

### History View
- Searchable list of all past commands
- Filter by: action type, date range, site
- Click to re-execute a command
- Pagination

### Browse View
- Embedded Chromium browser (BrowserView)
- Address bar (read-only, shows current URL)
- Back/forward buttons
- Opens when user says "open [site]" in command mode

### Settings View
- Microphone mode (Always / Push-to-talk / Off)
- Wake word (display only for now)
- Theme (dark only for v1)
- API connection status
- Ollama status
- Account management (email, linked devices)
- Logout / Delete account

---

## Chat → Ollama Flow

### New API Endpoint

```
POST /api/v1/chat
{
  "message": "What did I listen to yesterday?",
  "conversation_id": "uuid" (optional, for context)
}

Response:
{
  "reply": "Yesterday you played:\n• lofi beats on YouTube (9:14am)\n• jazz playlist on Spotify (2:30pm)",
  "sources": [
    { "id": "cmd-uuid", "transcript": "open youtube play lofi beats", "executed_at": "..." }
  ]
}
```

### Internal Flow (Go API)

```
1. Receive chat message
2. Send to Ollama with system prompt + user's recent history context
3. Ollama returns structured query OR direct answer
4. If structured query → execute against Postgres → format results
5. If direct answer → return as-is
6. Store conversation in DB for context
```

### Ollama System Prompt

```
You are Volo, a voice assistant's history helper. The user will ask questions
about their browsing and search history. You have access to their command history.

Given the user's question and their recent history (provided below), answer
naturally and concisely. If you need to search their history, output a JSON
query block like:

```json
{"action": "open-and-play", "date_from": "2026-05-31", "query_contains": "lofi"}
```

Available fields for queries:
- action: search, navigate, open-and-search, open-and-play, browser-control
- target: youtube, spotify, github, etc.
- query_contains: text search on the query field
- date_from, date_to: ISO date strings
- limit: max results (default 10)

Recent history (last 24h):
{INJECTED_HISTORY}
```

---

## New Database Tables

```sql
-- Chat conversations
CREATE TABLE conversations (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  title TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Chat messages
CREATE TABLE messages (
  id UUID PRIMARY KEY,
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  role TEXT NOT NULL,  -- 'user' or 'assistant'
  content TEXT NOT NULL,
  metadata JSONB DEFAULT '{}',  -- sources, query used, etc.
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at);
```

---

## Tech Stack

| Layer | Choice |
|-------|--------|
| Framework | Electron 30+ |
| UI | React 19 + TypeScript |
| Styling | Tailwind CSS |
| State | Zustand |
| Chat UI | assistant-ui components (or custom) |
| Voice | Web Speech API (renderer process) |
| Bundler | Vite (via electron-vite or electron-forge) |
| Packaging | electron-builder |
| LLM | Ollama (Gemma 2B) running locally |

---

## Docker Compose Addition

```yaml
ollama:
  image: ollama/ollama
  ports:
    - "11434:11434"
  volumes:
    - ollama_data:/root/.ollama
  # Pull model on first run:
  # docker exec ollama ollama pull gemma2:2b
```

---

## Implementation Order

1. Scaffold Electron + React + Vite
2. Main process: window, tray, global hotkey
3. Renderer: sidebar layout, routing between views
4. Chat view: message list, input, voice button
5. Wire to Volo API (same client as extension)
6. Add Ollama to Docker Compose
7. Implement /api/v1/chat endpoint in Go
8. Wire chat view → /api/v1/chat
9. History view (list + search + filters)
10. Browse view (embedded BrowserView)
11. Settings view
12. Packaging + auto-update

---

## Out of Scope (v1)

- Light theme
- Multiple conversations (single thread for now)
- File attachments
- Voice output (Volo doesn't speak back)
- Custom wake word
- Plugin system
