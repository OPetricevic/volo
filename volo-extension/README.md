# Volo Extension

Chromium browser extension for voice-activated browsing. Say "Hey Volo" followed by a command to search, navigate, open sites, and control your browser.

## Setup

```bash
npm install
```

## Build

```bash
npm run build        # Production build → dist/
npm run dev          # Watch mode (rebuilds on change)
```

## Load in Chrome

1. Run `npm run build`
2. Open `chrome://extensions`
3. Enable "Developer mode" (top right toggle)
4. Click "Load unpacked"
5. Select the `dist/` folder
6. Grant microphone permission when prompted

## Usage

- **Always listening mode:** Extension listens for "Hey Volo" continuously
- **Push-to-talk:** Click the extension icon or press `Ctrl+Shift+V`
- **Off:** Disable voice in the popup settings

### Supported Commands

| You say | What happens |
|---------|-------------|
| "Hey Volo, search [query]" | Google search |
| "Hey Volo, open [site]" | Navigate to site |
| "Hey Volo, open YouTube [query]" | YouTube search |
| "Hey Volo, open YouTube play [query]" | YouTube search + autoplay |
| "Hey Volo, go back" | Browser back |
| "Hey Volo, close tab" | Close current tab |
| "Hey Volo, new tab" | Open new tab |
| "Hey Volo, scroll down/up" | Scroll page |

## Testing

```bash
npm test             # Run all tests once
npm run test:watch   # Watch mode
```

## Project Structure

```
src/
├── background/       Service worker (command execution, badge, messaging)
├── content/          Content script (voice recognition, wake word, overlay)
├── popup/            Popup UI (React — status, mic mode toggle)
├── shared/           Shared logic (parser, types, storage)
├── styles/           Global CSS (Tailwind)
└── types/            TypeScript declarations (Web Speech API)
```

## Tech Stack

- React 19 + TypeScript
- Vite (bundler)
- Tailwind CSS
- Zustand (state)
- Web Speech API
- Manifest V3
