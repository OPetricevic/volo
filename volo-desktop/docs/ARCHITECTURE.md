# Desktop App Architecture

## Overview

Electron app with three process layers:

```
┌─────────────────────────────────────────────────┐
│ Main Process (Node.js)                          │
│                                                  │
│ • Window management (frameless, custom titlebar) │
│ • Ollama lifecycle (start/stop/model check)      │
│ • Global hotkey (Ctrl+Shift+V)                   │
│ • System tray                                    │
│ • IPC handler for renderer ↔ main communication  │
└────────────────────┬────────────────────────────┘
                     │ IPC (contextBridge)
┌────────────────────┴────────────────────────────┐
│ Preload (index.mjs)                             │
│                                                  │
│ Exposes window.volo API:                         │
│ • minimize(), maximize(), close()                │
│ • getOllamaStatus(), installOllama()             │
│ • onActivateVoice()                              │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────┐
│ Renderer (React + Tailwind)                     │
│                                                  │
│ Views:                                           │
│ • ChatView — AI chat with Ollama (bubble layout) │
│ • HistoryView — browsing history viewer          │
│ • MacrosView — create/edit custom voice commands │
│ • SettingsView — mic mode, connections, account  │
│                                                  │
│ Sidebar: Chat / History / Macros / Settings      │
└─────────────────────────────────────────────────┘
```

## Ollama Integration

Ollama is bundled via the NSIS installer. On app startup:

1. `OllamaManager.initialize()` checks if Ollama is running at `localhost:11434`
2. If not running but installed → starts `ollama.exe serve`
3. If model missing → pulls `gemma2:2b` (~1.6GB first-time download)
4. If Ollama not found → logs warning, chat disabled (voice commands still work via extension)

**Binary location:** `%LOCALAPPDATA%\Programs\Ollama\ollama.exe` (standard Ollama install path)

## Macros

Macros are custom voice commands managed in the desktop app and executed by the browser extension.

### Flow

```
User creates macro in MacrosView
  → POST /api/v1/macros (saved to Postgres)
  → Extension fetches on next cache refresh (5min TTL)
  → User says "Hey Volo, [trigger phrase]"
  → Extension matches trigger → opens all URLs
```

### Limits

- Max 10 macros per user
- Max 3 actions (URLs) per macro
- Requires internet (synced via API)
- Only logged-in users can create macros

### Offline behavior

Shows a yellow warning banner: "Offline — showing demo macro."
A pre-configured "Morning Routine" macro is shown as an example.
Create/edit/delete buttons are disabled when offline.

## Chat (Ollama AI)

Each message is **stateless** — no conversation memory. The prompt includes:
- System instructions (Volo is a history helper)
- Recent command history from the API (`GET /history/context`)
- The user's current question

Gemma 2B responds based only on what's in the prompt. No multi-turn context.

**Context budget:** ~4000-6000 tokens for history (8192 window total).
**Response cap:** 256 tokens (`num_predict` option).

## Color System

CSS custom properties in `globals.css`:

| Variable | Value | Usage |
|----------|-------|-------|
| `--bg` | `#0a0a0c` | App background |
| `--sidebar` | `#101112` | Sidebar background |
| `--chat` | `#121314` | Chat/main area background |
| `--input` | `rgba(255,255,255,0.03)` | Input fields, cards |
| `--border` | `rgba(255,255,255,0.06)` | Borders, dividers |
| `--text` | `#d0d6e0` | Primary text |
| `--text-muted` | `#9c9da1` | Secondary text |
| `--text-faint` | `#62666d` | Labels, placeholders |
| `--accent` | `#0ea5e9` | Interactive elements |
| `--success` | `#22c55e` | Online/connected states |
| `--warning` | `#f59e0b` | Warnings |
| `--danger` | `#ef4444` | Errors, destructive |

Font: Inter (loaded from Google Fonts).

## Packaging

- **Build:** `npx electron-vite build`
- **Package:** `npm run package:win` → produces `release/Volo-Setup-0.1.0.exe`
- **Installer:** NSIS (downloads Ollama + pulls model during install)
- **Target:** Windows 10+ only
- **System requirements:** 16GB RAM recommended, ~3GB disk for AI model
