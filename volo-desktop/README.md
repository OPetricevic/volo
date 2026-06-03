# Volo Desktop

Electron desktop app with Discord-style dark UI and Claude/Codex-style chat layout. Chat with your browsing history using a local Gemma model via Ollama.

## Setup

```bash
npm install
```

## Development

```bash
npm run dev
```

Opens the Electron app with hot-reload.

## Build

```bash
npm run build
```

Outputs to `out/` (main, preload, renderer).

## Package (create installer)

```bash
npm run package:win
```

Creates `.exe` installer in `release/`.

## Test

```bash
npm test
```

20 tests covering config, prompt building, and Ollama logic.

## Architecture

- **Main process** (`src/main/`) — Window management, tray, global hotkey, Ollama manager, auto-updater
- **Preload** (`src/preload/`) — Safe IPC bridge between main and renderer
- **Renderer** (`src/renderer/`) — React UI (chat, history, settings)

## UI Design

- **Colors:** Discord dark palette (`#1e1f22`, `#2b2d31`, `#313338`) with Tailwind blue accent (`#3b82f6`)
- **Layout:** Claude/Codex style — clean sidebar with conversations, centered chat, spacious messages
- **Not Discord-like:** No server icons, no channels, no status. Just a clean chat app.

## Key Features

- Chat with Volo AI about your search history (local Ollama + Gemma 2B)
- Voice commands (same as extension — "Hey Volo, open YouTube")
- Command history view with search/filter
- Settings (mic mode, API connection, Ollama status)
- System tray (minimize to tray, global hotkey Ctrl+Shift+V)
- Auto-update via GitHub Releases

## Config

API URL defaults to production. Override for development:

```bash
# .env
VITE_API_URL=http://localhost:8080
```

## Ollama (AI Chat)

Optional. If installed locally:
- Desktop calls Ollama directly (`localhost:11434`)
- API only provides history context (lightweight)
- Server stays lean (no AI compute on VM)

Install: checkbox during Volo setup, or manually install Ollama + `ollama pull gemma2:2b`
