# TODO — Volo

## In Progress

- [ ] Desktop: Settings view needs update to match new sidebar style

---

## UI / Visual

- [ ] Desktop: Settings view needs update to match new sidebar style
- [ ] Desktop: History view needs update to match new style
- [ ] Landing page: responsive polish (mobile nav, mockup scaling)
- [ ] Landing page: real copy review (headlines, descriptions)
- [ ] Icons: real app icon design (extension + desktop) — replace placeholders

---

## Functional

(all functional items complete)

---

## ML / AI

- [ ] Run ML training (Docker: prepare_dataset → train → export_onnx)
- [ ] Copy ONNX file to volo-api/models/intent.onnx + vocab + config
- [ ] Deploy with ONNX Runtime on Linux (replace onnx_stub.go with real implementation)

### ML Steps (requires Docker running)

```bash
# 1. Start the ML container
docker compose -f docker-compose.dev.yml up -d model

# 2. Prepare dataset (downloads SNIPS, generates custom examples, augments)
docker compose -f docker-compose.dev.yml exec model python scripts/prepare_dataset.py

# 3. Train DistilBERT (~10 min on CPU)
docker compose -f docker-compose.dev.yml exec model python scripts/train.py

# 4. Export to ONNX + quantize
docker compose -f docker-compose.dev.yml exec model python scripts/export_onnx.py

# 5. Copy model files to API
cp volo-model/models/volo-intent-quantized.onnx volo-api/models/intent.onnx
cp -r volo-model/models/tokenizer volo-api/models/tokenizer
cp volo-model/models/model-config.json volo-api/models/model-config.json

# 6. On Linux deployment: replace volo-api/internal/intent/onnx_stub.go
#    with a real implementation using github.com/yalue/onnxruntime_go
#    (stub exists because onnxruntime doesn't compile on Windows)

# 7. Set env var and restart API
# VOLO_MODEL_PATH=./models
# docker compose restart volo-api
```

---

## Infrastructure

- [ ] Deploy to VM (see DEPLOYMENT.md for full guide)

---

## Configuration (manual, one-time setup)

- [ ] Google OAuth: Create project in Google Cloud Console → enable Identity API → create OAuth 2.0 credentials → set VOLO_GOOGLE_CLIENT_ID + VOLO_GOOGLE_CLIENT_SECRET in .env
- [ ] Sentry: Create account at sentry.io → create Go project → copy DSN → set VOLO_SENTRY_DSN in .env
- [ ] Domain: Point api.volo.yourdomain.com + volo.yourdomain.com to VM IP (A records)
- [ ] JWT secret: Generate with `openssl rand -hex 32` → set VOLO_JWT_SECRET in .env
- [ ] Redis password: Generate with `openssl rand -hex 16` → set VOLO_REDIS_PASSWORD in .env
- [ ] Chrome Web Store: Create developer account ($5 one-time) for extension publishing
- [ ] GitHub Releases: Set up repo for desktop auto-updater to pull from
- [ ] Prometheus server: Add to VM docker-compose (see DEPLOYMENT.md)
- [ ] Grafana: Add to VM docker-compose, create dashboards (see DEPLOYMENT.md)
- [ ] Caddy + TLS: Configure reverse proxy on VM (see DEPLOYMENT.md Caddyfile)
- [ ] Postgres backups: Set up cron job on VM (see DEPLOYMENT.md backups section)

---

## Testing

- [ ] End-to-end test: speak → extension → API → Postgres → verify
- [ ] Desktop app: manual test with Ollama running
- [ ] Extension: manual test in Chrome (load, speak, verify)
- [ ] Load test API (100 concurrent users simulation)

---

## Packaging / Release

- [ ] Desktop: NSIS installer with optional Ollama checkbox
- [ ] Desktop: code signing (Windows)
- [ ] Extension: publish to Chrome Web Store (developer account needed)
- [ ] Landing page: deploy to VM or Vercel

---

## Documentation

- [ ] Create `volo-docs` repo for visual HTML documentation (GitHub Pages)
- [ ] Redesign docs site with C4 model structure + ByteByteGo visual style
- [ ] Host on GitHub Pages: https://opetricevic.github.io/volo-docs
- [ ] Link from main volo README to the docs site

### Docs Site Redesign Plan

**Style references:**
- ByteByteGo (bytebytego.com/guides) — clean SVG diagrams, numbered flows, color-coded
- C4 Model (c4model.com) — zoom levels: Context → Container → Component
- Arc42 template — structured sections (context, decisions, quality, deployment)
- IcePanel (icepanel.io) — interactive clickable architecture diagrams
- Example repo: github.com/bitsmuggler/arc42-c4-software-architecture-documentation-example

**Pages to create/redesign:**
1. System Context (Level 1) — Volo + users + external systems (Google, Chrome, Ollama)
2. Container Diagram (Level 2) — Extension, Desktop, API, Model, DB, Redis
3. Voice Command Flow — mic → wake word → parse → execute → history (numbered steps)
4. Auth Architecture — device registration → JWT → session → revocation (flow diagram)
5. AI Pipeline — rule-based → confidence check → DistilBERT → response (decision tree)
6. Data Architecture — Postgres schema relationships, Redis caching strategy
7. Deployment — VM layout, Docker containers, Caddy, monitoring (infrastructure diagram)
8. Decisions — 16 architecture decisions (already exists, keep as-is)

**Each page should have:**
- One clear SVG diagram at the top (not ASCII, not text boxes — proper vector graphics)
- Numbered steps explaining the flow
- "Why this design" section with tradeoffs
- Color-coded components (blue=extension, green=API, purple=model, orange=desktop)
- Clickable/interactive where possible (hover for details, tabs for alternatives)

**Current state:** Pages exist with content but use text-based diagrams. Need visual upgrade.

---

## Done ✓

- [x] Phase 1 — Extension (built, tested, 50 tests)
- [x] Phase 2 — Go API (built, tested, 105+ tests)
- [x] Phase 3 — ML pipeline scripts
- [x] Phase 4 — Desktop app scaffold
- [x] Phase 5 — Landing page
- [x] Extension ↔ API wiring
- [x] Sentry + Prometheus + validation middleware
- [x] DECISIONS.md + all documentation
- [x] Desktop UI: Discord palette + Codex layout
- [x] Landing page: animated demo, interactive send
- [x] Extension UI: updated to Discord palette
- [x] Extension retry queue (9 tests)
- [x] Google OAuth implementation (handler + service)
- [x] Platform-grade tests (concurrency, auth validation, edge cases)
- [x] Suggestions endpoint (frequency * recency scoring, top 10, 5 tests)
- [x] Settings sync (GET/PUT /settings, JSONB storage, partial updates)
- [x] Extension onboarding (welcome screen + mic permission chooser on first launch)
- [x] Go ONNX integration (hybrid parser, model loader, tokenizer, stub for non-Linux, 5 tests)
- [x] Prometheus metrics fully wired (RecordCommand, RecordRegistration, RecordRateLimitHit, RecordOllamaRequest)
