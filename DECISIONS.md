# Architecture Decisions

A log of every significant technical decision made in this project, with reasoning. Reference this when picking the project back up to understand *why* things are the way they are.

---

## 1. Hybrid Intent Parser (Rule-Based + ML)

**Decision:** Use a rule-based parser as the fast path (~1ms), fall back to DistilBERT ONNX (~8ms) when confidence < 0.80.

**Why:**
- Rule-based handles 80%+ of commands perfectly (they follow predictable patterns)
- ML model only fires on ambiguous/novel inputs (saves compute)
- User gets instant response for common commands (no network round-trip to model)
- If the model is unavailable, everything still works

**Alternatives considered:**
- LLM for everything (too slow, too expensive, overkill for "open youtube")
- Rule-based only (can't handle "find that thing from yesterday")
- ML only (unnecessary latency for simple commands)

---

## 2. DistilBERT (not GPT, not TinyBERT)

**Decision:** Fine-tune DistilBERT (67M params) for intent classification.

**Why:**
- Small enough for CPU inference in ~8ms
- Large enough for 92-96% accuracy on 5 intents
- Well-documented, MIT licensed, good for thesis citations
- ONNX export is straightforward
- No GPU required anywhere in the pipeline

**Alternatives considered:**
- GPT-4o-mini API (costs money, needs internet, overkill)
- TinyBERT (14M params, accuracy drops too much)
- Phi-3 (3.8B params, needs 4GB+ RAM, too heavy for classification)

---

## 3. Ollama + Gemma for Chat (not the same model)

**Decision:** Use Ollama running Gemma 2B locally for the chat/history query feature. Separate from the intent classifier.

**Why:**
- Chat requires *generating* natural language responses (classification can't do this)
- Gemma 2B is small enough to run on CPU (~2.3GB)
- Ollama manages the model lifecycle (download, serve, memory)
- Runs locally — no API costs, no internet dependency, privacy-first
- Optional component — app works without it (just no chat)

**Why not use Gemma for intent classification too?**
- Overkill. DistilBERT classifies in 8ms. Gemma would take 500ms+ for the same task.
- Different problems need different tools.

---

## 4. PostgreSQL (not SQLite)

**Decision:** Use Postgres for all persistent data, even in development.

**Why:**
- Treating this as production-grade from day one (designed for 1000+ users)
- Postgres handles concurrent writes properly (SQLite doesn't)
- JSONB columns for flexible data (hour_weights, metadata)
- Proper indexing, full-text search if needed later
- Same database in dev and prod = no "works on my machine" issues

**Alternatives considered:**
- SQLite (simpler, but single-writer, no concurrent access)
- MongoDB (overkill, no relational integrity)

---

## 5. Redis for Rate Limiting + Caching (not in-memory)

**Decision:** Use Redis for rate limiting, session caching, and lookup table caching.

**Why:**
- Survives API restarts (in-memory counters reset on deploy)
- Sliding window rate limiting needs sorted sets (Redis is perfect for this)
- Session validation cache reduces Postgres queries by ~90%
- If we ever scale to multiple API instances, Redis is shared state

**What Redis stores:**
- Rate limit counters (sorted sets, per user, sliding window)
- Session cache (`session:<token_hash>` → `user_id`, TTL = token expiry)
- Lookup tables (sites, action_types — loaded on boot)

---

## 6. UUIDs Generated in Go (not Postgres)

**Decision:** All UUIDs are generated in the Go service layer using `google/uuid`, not via `gen_random_uuid()` in Postgres.

**Why:**
- No dependency on `pgcrypto` extension
- Testable (can control/mock IDs in tests)
- Consistent behavior regardless of database
- The service layer owns identity creation, not the database

---

## 7. Passwords in Separate Table (not in users)

**Decision:** Auth credentials live in a `credentials` table, not in `users`.

**Why:**
- Security separation — `users` table has NO secrets
- One user can have multiple auth methods (password + Google + GitHub)
- Adding a new OAuth provider = just another row, no schema change
- If credentials table is compromised, user data is still safe
- Clean separation of concerns

---

## 8. Offline-First Extension

**Decision:** The extension always works locally, even without the backend. API is an enhancement, not a dependency.

**Why:**
- Voice commands should be instant (can't wait for network)
- If the API is down, users shouldn't notice for basic commands
- History/learning is nice-to-have, not critical path
- Reduces perceived latency to near-zero

**How it works:**
1. Command detected → parse locally → execute immediately (~1ms)
2. In background: send to API for history storage (fire-and-forget)
3. If API is unreachable, command still worked. History is lost for that command.

---

## 9. Sentry Cloud (not self-hosted)

**Decision:** Use Sentry's cloud service for error tracking.

**Why:**
- Free tier (5K errors/month) is more than enough
- Zero maintenance — no extra container to manage
- Instant setup (just add DSN env var)
- Can switch to self-hosted later if needed (same SDK)

**When to reconsider:** If error volume exceeds free tier, or if data residency requirements change.

---

## 10. Electron (not Tauri)

**Decision:** Use Electron for the desktop app.

**Why:**
- We need an embedded Chromium browser (BrowserView) — Electron has this built-in
- Shares React code with the extension (same UI framework)
- Proven at scale (Discord, Slack, VS Code all use it)
- Web Speech API works in Electron's Chromium renderer
- Tauri uses system webview which doesn't support Web Speech API consistently

**Tradeoff:** Larger binary (~80MB vs Tauri's ~5MB). Acceptable for a desktop app.

---

## 11. Ollama as Optional Installer Component

**Decision:** The NSIS installer has a single checkbox "Enable Volo AI" that downloads Ollama + Gemma during install.

**Why:**
- Volo Desktop is ~80MB. Adding Ollama + model makes it ~2.5GB.
- Most users want the voice commands (small download). AI chat is a power feature.
- Users who don't check the box get a fully functional app (just no chat)
- Can always enable later from Settings
- No separate Ollama install step — it's bundled inside Volo's directory

---

## 12. Single Docker Compose for Dev Environment

**Decision:** One `docker-compose.dev.yml` at the project root spins up everything: Postgres, Redis, Go API (hot-reload), Python ML environment (Jupyter), Ollama.

**Why:**
- New machine? `docker compose up -d` and you're working in 2 minutes
- No "install Go, install Python, install Postgres" instructions
- Everyone gets the same versions of everything
- Works on Windows, Mac, Linux identically
- The ML training runs in a container too — no Python version conflicts

---

## 13. SNIPS Dataset + Custom Examples

**Decision:** Use the SNIPS NLU benchmark (MIT licensed) as base training data, supplemented with ~200 custom browser command examples.

**Why:**
- SNIPS has relevant intents (PlayMusic, SearchCreativeWork) that map to our labels
- Well-studied in academia — good for thesis citations
- MIT licensed — no legal issues
- Custom examples cover browser-specific commands SNIPS doesn't have
- Augmentation with filler words ("um", "uh") simulates real speech recognition output

---

## 14. Error Response Format (Dual Errors)

**Decision:** Every error response includes both a user-friendly `message` AND a detailed `internal_error` with the full error chain.

```json
{
  "error": {
    "code": "COMMAND_PARSE_FAILED",
    "message": "We couldn't process that command. Try again.",
    "internal_error": "handler.ProcessCommand → service.Command.Process → intent.Classify: tokenizer failed"
  }
}
```

**Why:**
- Frontend shows `message` to the user (friendly, actionable)
- Developer reads `internal_error` for debugging (exact failure point)
- No need to dig through server logs for most issues
- The error chain tells you exactly which layer failed
- Audit log stores the full chain for historical debugging

---

## 15. Soft Deletes Everywhere

**Decision:** All deletable entities use `deleted_at TIMESTAMPTZ` instead of hard DELETE.

**Why:**
- Audit trail — can always see what existed
- Recoverable — "oops I deleted my account" is fixable
- No FK violations — related data doesn't break
- GDPR compliance — can prove data was "deleted" (set timestamp) while keeping audit
- Queries filter with `WHERE deleted_at IS NULL` (indexed)

---

## 16. JWT with Session Table (not stateless)

**Decision:** JWTs are stored in a `sessions` table (token hash). Validation checks the table.

**Why:**
- **Revocable:** Logout immediately invalidates the token (stateless JWTs can't do this)
- **Visible:** Can see all active sessions, revoke specific devices
- **Secure:** If a token is stolen, you can kill it without rotating the secret
- Redis caches the validation (fast path), Postgres is source of truth (revocation)

**Tradeoff:** Slightly slower than pure stateless JWT (Redis lookup). But revocability is non-negotiable for a real app.
