# Extension Testing

## Overview

41 tests covering the core logic. Uses Vitest as the test runner.

## Running Tests

```bash
npm test             # Single run
npm run test:watch   # Watch mode (re-runs on file change)
```

## Test Coverage

### Wake Word Detection (8 tests)

| Test | What it verifies |
|------|-----------------|
| Detects at start | "hey volo search something" → true |
| Detects in middle | "um hey volo open youtube" → true |
| Case insensitive | "HEY VOLO" → true |
| Absent | "hey siri open youtube" → false |
| Partial match | "hey vol" → false |
| Extracts command | "hey volo search react" → "search react" |
| Handles whitespace | "  hey volo   open youtube  " → "open youtube" |
| Nothing after wake | "hey volo" → "" |

### Command Parsing (27 tests)

**Search commands (5):** `search X`, `search for X`, `look up X`, `find X`, fallback

**Navigate commands (6):** `open youtube`, `go to github`, `goto reddit`, fuzzy "tube"→youtube, fuzzy "yt"→youtube, raw URLs

**Open + search (3):** `open youtube lofi`, `open spotify jazz`, `open github react`

**Open + play (2):** `open youtube play lofi hip hop`, `open spotify play jazz`

**Browser control (8):** go back, back, forward, close tab, close this tab, new tab, scroll down, scroll up

**Case insensitivity (2):** uppercase, mixed case

### URL Building (6 tests)

| Test | Input | Expected URL |
|------|-------|-------------|
| Google search | search "react hooks" | google.com/search?q=react%20hooks |
| YouTube nav | navigate "youtube" | youtube.com |
| YouTube search | open-and-search youtube "lofi" | youtube.com/results?search_query=lofi%20beats |
| Spotify search | open-and-search spotify "jazz" | open.spotify.com/search/jazz |
| Browser control | browser-control "back" | null |
| Raw URL | navigate "https://example.com" | https://example.com |

## What's NOT Tested (requires browser)

- Web Speech API integration (needs real microphone)
- Chrome extension messaging (needs Chrome runtime)
- Popup React rendering (would need jsdom + more setup)
- Overlay DOM injection

These are tested manually by loading the extension in Chrome.
