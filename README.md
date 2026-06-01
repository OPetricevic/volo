<p align="center">
  <img src="https://img.shields.io/badge/version-0.1.0-blue" alt="Version">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
  <img src="https://img.shields.io/badge/platform-Windows-lightgrey" alt="Platform">
</p>

<h1 align="center">Volo</h1>

<p align="center">
  <strong>Voice-first browser assistant.</strong><br>
  Say "Hey Volo" to search, navigate, and control your browser with your voice.<br>
  Ask about your history in natural language. It learns your patterns over time.
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#development">Development</a> •
  <a href="#documentation">Documentation</a>
</p>

---

## Quick Start

```bash
# Clone
git clone https://github.com/yourname/volo.git
cd volo

# Start backend (Postgres, Redis, API)
docker compose -f docker-compose.dev.yml up -d

# Build extension
cd volo-extension && npm install && npm run build
# Load dist/ folder in chrome://extensions (Developer mode)

# Run desktop app
cd volo-desktop && npm install && npm run dev

# Landing page
cd volo-landing && npm install && npm run dev
```

## Features

- **Voice Commands** — Say "Hey Volo, search React hooks" or "open YouTube play lofi"
- **AI History Chat** — Ask "What did I listen to yesterday?" and get answers from your own data
- **Instant Execution** — Commands execute in <10ms locally, no round-trip needed
- **Pattern Learning** — Learns your habits, surfaces frequent commands faster
- **Cross-Platform** — Browser extension (Chrome/Edge/Opera) + Desktop app (Windows)
- **Privacy First** — Voice processing and AI run locally on your machine
- **Offline-First** — Works without internet, syncs when connected

## Architecture

```
                    ┌─────────────────────────┐
                    │         User            │
                    │  "Hey Volo, open YouTube │
                    │   play lofi"            │
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
          ┌─────────┴──────────┐   ┌─────────┴──────────┐
          │ Browser Extension  │   │   Desktop App       │
          │ (Chrome/Edge/Opera)│   │   (Electron)        │
          │                    │   │                     │
          │ • Voice recognition│   │ • Voice + AI chat   │
          │ • Wake word detect │   │ • Local Ollama      │
          │ • Command execution│   │ • Embedded browser  │
          └─────────┬──────────┘   └─────────┬──────────┘
                    │                         │
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │        Go API           │
                    │                         │
                    │ • Auth (JWT + Google)    │
                    │ • Command history        │
                    │ • Pattern learning       │
                    │ • Suggestions            │
                    │ • Postgres + Redis       │
                    └─────────────────────────┘
```

## Project Structure

```
volo/
├── volo-extension/     Browser extension (React, Manifest V3, Web Speech API)
├── volo-desktop/       Desktop app (Electron, React)
├── volo-api/           Backend (Go, Chi, Postgres, Redis)
├── volo-model/         ML training pipeline (Python, DistilBERT, ONNX)
├── volo-landing/       Marketing site (React, Vite, Tailwind)
├── docker-compose.dev.yml   Full dev environment
├── SPEC.md             Project specification
├── DECISIONS.md        Architecture decisions
├── TODO.md             Current task list
└── DOCS.md             Documentation index
```

## Development

### Prerequisites

- Node.js 20+
- Go 1.22+
- Docker (for Postgres, Redis, ML training)

### Running Tests

```bash
# Extension (50 tests)
cd volo-extension && npm test

# API (105+ tests)
cd volo-api && go test -short ./...

# Desktop (20 tests)
cd volo-desktop && npm test
```

### Full Stack (Docker)

```bash
docker compose -f docker-compose.dev.yml up -d
```

This starts: Postgres (5432), Redis (6379), Go API with hot-reload (8080), Jupyter for ML (8888), Ollama (11434).

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Extension | React, TypeScript, Manifest V3, Web Speech API, Vite |
| Desktop | Electron, React, TypeScript, Tailwind |
| API | Go, Chi, PostgreSQL, Redis, JWT |
| ML | Python, PyTorch, HuggingFace, DistilBERT, ONNX |
| Landing | React, TypeScript, Vite, Tailwind |
| Infra | Docker, Caddy, Prometheus, Grafana, Sentry |

## Documentation

| Document | Description |
|----------|-------------|
| [SPEC.md](./SPEC.md) | Full project specification |
| [DECISIONS.md](./DECISIONS.md) | 16 architecture decisions with reasoning |
| [DOCS.md](./DOCS.md) | Documentation index (links to all docs) |
| [TODO.md](./TODO.md) | Current task list |
| [DOCS.html](./DOCS.html) | Visual documentation (open in browser) |
| [SPEC.html](./SPEC.html) | Visual spec (open in browser) |

## Contributing

This is a bachelor thesis project. Contributions welcome after initial release.

## License

MIT
