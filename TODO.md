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
- [ ] Convert all MD docs to interactive HTML pages (same style as the deleted SPEC.html/DOCS.html)
- [ ] Host on GitHub Pages: https://opetricevic.github.io/volo-docs
- [ ] Link from main volo README to the docs site
- [ ] Pages to create: Architecture, API Reference, Database Schema, Auth Flows, Intent Parser, Deployment, Observability
- [ ] Make them interactive (tabs, collapsible sections, SVG diagrams, code highlighting)

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
