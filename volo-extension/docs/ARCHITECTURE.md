# Extension Architecture

## Overview

The extension has three execution contexts that communicate via Chrome's messaging API:

```
┌─────────────────────────────────────────────────────────┐
│                    Browser                                │
│                                                          │
│  ┌──────────────┐   messages   ┌────────────────────┐   │
│  │ Service Worker│ ←──────────→ │ Content Script      │   │
│  │ (background)  │              │ (per tab)           │   │
│  │               │              │                     │   │
│  │ • Command     │              │ • Web Speech API    │   │
│  │   execution   │              │ • Wake word detect  │   │
│  │ • Tab control │              │ • Overlay UI        │   │
│  │ • Badge state │              │ • Scroll commands   │   │
│  └───────┬───────┘              └─────────────────────┘   │
│          │ messages                                       │
│  ┌───────┴───────┐                                       │
│  │ Popup (React) │                                       │
│  │               │                                       │
│  │ • Status      │                                       │
│  │ • Mic toggle  │                                       │
│  │ • Quick help  │                                       │
│  └───────────────┘                                       │
└─────────────────────────────────────────────────────────┘
```

## Content Script (`src/content/content.ts`)

Injected on every page. Responsibilities:

1. **Voice recognition** — Creates a `SpeechRecognition` instance, listens continuously
2. **Wake word detection** — Scans interim transcripts for "hey volo"
3. **Command capture** — After wake word, buffers speech until 2s silence or final result
4. **Overlay UI** — Shows a floating "Listening..." pill when wake word is detected
5. **Local commands** — Handles scroll up/down directly (no round-trip needed)

### Lifecycle

```
Page loads → init() → check mic mode from storage
  → if "always": startListening()
  → if "once": wait for ACTIVATE_LISTENING message
  → if "off": do nothing

Listening loop:
  recognition.onresult → check for wake word
    → found: set wakeWordDetected=true, show overlay, start silence timer
    → accumulate command buffer
    → on final result OR 2s silence: processCommand()
      → send COMMAND_DETECTED to background
      → reset state, hide overlay

  recognition.onend → if "always" mode, restart after 300ms
```

## Service Worker (`src/background/service-worker.ts`)

Persistent background process. Responsibilities:

1. **Message routing** — Receives commands from content script, dispatches actions
2. **Command parsing** — Runs the rule-based parser on transcripts
3. **Tab control** — Opens URLs, navigates back/forward, closes/creates tabs
4. **Badge management** — Updates icon badge color based on voice state
5. **Keyboard shortcut** — Listens for Ctrl+Shift+V to activate voice

### Message Types

| Message | From | To | Purpose |
|---------|------|-----|---------|
| `VOICE_STATE_CHANGED` | Content | Background | Update badge |
| `COMMAND_DETECTED` | Content | Background | Process command |
| `GET_STATE` | Popup | Background | Get current state |
| `SET_MIC_MODE` | Popup | Background → All tabs | Change listening mode |
| `ACTIVATE_LISTENING` | Popup/Shortcut | Background → All tabs | Start voice |
| `EXECUTE_COMMAND` | Background | Content | Scroll commands |

## Popup (`src/popup/`)

React app rendered in the extension popup (320px wide). Shows:

- Current voice state (idle/listening/processing)
- Mic mode selector (Always / Push-to-talk / Off)
- Activate button (in push-to-talk mode)
- Quick help with example commands

## Parser (`src/shared/parser.ts`)

Rule-based intent classifier. Runs client-side for zero-latency response.

Priority order:
1. Browser control (exact match: "go back", "close tab", etc.)
2. Open + play (`open [site] play [query]`)
3. Open + search (`open [site] [query]`)
4. Navigate (`open [site]` or `go to [site]`)
5. Search (`search [query]`, `search for [query]`, `find [query]`)
6. Fallback (treat entire text as Google search, confidence 0.6)

Includes fuzzy site matching: "tube" → youtube, "yt" → youtube, "gh" → github.
