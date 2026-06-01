# Observability & Error Tracking

## Why This Exists

Without observability, if the app breaks in production you have no way of knowing until a user complains. This is unacceptable for a production-grade system. We need to answer three questions at any time:

1. **Is it working?** (health checks, uptime)
2. **How well is it working?** (latency, error rate, throughput)
3. **What broke and why?** (error traces, stack traces, context)

---

## Architecture Decision

We chose a three-layer approach:

| Layer | Tool | Why |
|-------|------|-----|
| Error tracking | Sentry (cloud) | Captures crashes with full stack traces, groups duplicates, alerts on new issues. Free tier (5K errors/month) is enough for a thesis. |
| Metrics | Prometheus + Grafana | Industry standard. Prometheus scrapes numbers, Grafana visualizes them. Both run as Docker containers on the VM. |
| Structured logging | Go `slog` (JSON) | Every request gets a JSON log line with request_id, method, path, status, duration. Queryable if needed. |

### Why not just logs?

Logs tell you what happened *after* you go looking. Metrics tell you something is wrong *right now*. Sentry tells you *exactly what broke* with the stack trace. You need all three.

### Why Sentry Cloud vs Self-Hosted?

- Self-hosted Sentry needs ~2GB RAM and is complex to maintain
- Cloud free tier gives 5K errors/month — more than enough
- Zero maintenance, just add the DSN env var
- Can switch to self-hosted later if needed

---

## Sentry Integration

### How It Works

```
User makes request → API processes it → Something throws an error
                                              │
                                              ▼
                                    Sentry SDK captures:
                                    • Stack trace
                                    • Request context (method, path, headers)
                                    • User ID (from JWT)
                                    • Environment (dev/prod)
                                    • Release version
                                              │
                                              ▼
                                    Sends to Sentry cloud
                                              │
                                              ▼
                                    You get an email/Slack notification
```

### What Gets Captured

| Event | How |
|-------|-----|
| Unhandled panics | `SentryMiddleware()` + `Recoverer` catches panics, reports to Sentry, then returns 500 |
| Explicit errors | `CaptureError(err, tags)` called in service layer for important failures |
| Performance traces | 20% of requests get full timing breakdown (where time was spent) |

### Configuration

```bash
# Set in .env on the VM
VOLO_SENTRY_DSN=https://abc123@o123456.ingest.sentry.io/789
VOLO_ENVIRONMENT=production  # or "development"
```

If `VOLO_SENTRY_DSN` is empty, Sentry is completely disabled. No impact on performance or behavior.

### Setup Steps

1. Create account at [sentry.io](https://sentry.io)
2. Create a Go project
3. Copy the DSN
4. Set `VOLO_SENTRY_DSN` in your `.env`
5. Deploy — errors start appearing in Sentry dashboard

---

## Prometheus Metrics

### How It Works

```
Prometheus container (every 15s)
    │
    │  GET /metrics
    ▼
Volo API responds with current metric values
    │
    │  stored in Prometheus time-series DB
    ▼
Grafana queries Prometheus → renders dashboards
```

### Metrics Exposed

| Metric | Type | Labels | What it tells you |
|--------|------|--------|-------------------|
| `volo_http_requests_total` | Counter | method, path, status | Total requests. Alert if 5xx rate > 5%. |
| `volo_http_request_duration_seconds` | Histogram | method, path | Latency distribution. Alert if P95 > 2s. |
| `volo_commands_processed_total` | Counter | action | Which commands are used most. Product insight. |
| `volo_auth_registrations_total` | Counter | — | New users over time. Growth metric. |
| `volo_rate_limit_hits_total` | Counter | — | Abuse detection. Spike = someone hammering the API. |
| `volo_ollama_requests_total` | Counter | status | Ollama health. If errors spike, model is down. |

### Grafana Alerts (Recommended)

| Alert | Condition | Severity |
|-------|-----------|----------|
| API Down | `/health` returns non-200 for > 1 min | Critical |
| High Error Rate | 5xx rate > 5% for 5 min | High |
| High Latency | P95 > 2s for 5 min | Medium |
| Rate Limit Spike | > 100 hits in 5 min | Medium |
| Ollama Down | Error rate 100% for 2 min | Low (chat degrades gracefully) |

### Configuration

Prometheus config is in `prometheus/prometheus.yml`:
```yaml
scrape_configs:
  - job_name: 'volo-api'
    static_configs:
      - targets: ['volo-api:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

Grafana is accessible at `:3000` on the VM (not exposed publicly unless you add a subdomain).

---

## Input Validation

### Why

Without validation, the API is vulnerable to:
- **Memory exhaustion:** 100MB request body fills RAM
- **Storage abuse:** 10,000-character transcripts fill the database
- **Brute force:** Unlimited login attempts
- **Garbage data:** Empty strings stored as valid commands

### What's Validated

| Field | Constraint | Error Code |
|-------|-----------|-----------|
| Request body | Max 10KB | `BODY_TOO_LARGE` |
| Transcript | 1–500 chars, not whitespace-only | `EMPTY_TRANSCRIPT`, `TRANSCRIPT_TOO_LONG` |
| Chat message | 1–2000 chars | `EMPTY_MESSAGE`, `MESSAGE_TOO_LONG` |
| Email | Must contain @ and ., max 255 chars | `MISSING_EMAIL`, `INVALID_EMAIL`, `EMAIL_TOO_LONG` |
| Password | 8–128 chars | `WEAK_PASSWORD`, `PASSWORD_TOO_LONG` |

### Implementation

- `ValidateBody` middleware runs on ALL requests (global)
- Field validators are called in handlers before processing
- All return structured `ErrorResponse` with code + message
- 16 tests cover all boundary conditions

---

## Structured Logging

### Format

Every request produces a JSON log line:
```json
{
  "level": "INFO",
  "msg": "request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/command",
  "status": 200,
  "duration_ms": 12,
  "ip": "192.168.1.100"
}
```

Errors include the full chain:
```json
{
  "level": "ERROR",
  "msg": "command processing failed",
  "request_id": "...",
  "error_chain": "handler.ProcessCommand → service.Command.Process → repository.Command.Create: connection refused"
}
```

### Why JSON?

- Machine-parseable (can pipe to log aggregation tools later)
- Structured fields (filter by status, path, duration)
- Request ID links all logs for a single request together

---

## What Happens When Things Break

### Scenario: API crashes

1. Recoverer middleware catches the panic
2. Sentry captures the stack trace + request context
3. Returns 500 JSON to the user (not a blank page)
4. Prometheus records the 5xx
5. Grafana alert fires → you get notified
6. You check Sentry → see exactly which line panicked

### Scenario: Ollama is down

1. Chat endpoint calls Ollama → timeout after 30s
2. Returns error to user: "Volo AI couldn't process that"
3. `volo_ollama_requests_total{status="error"}` increments
4. Grafana shows Ollama error rate spike
5. Commands still work (they don't need Ollama)

### Scenario: Database is full

1. INSERT fails → repository returns error
2. Error wrapping shows: `repository.Command.Create: disk full`
3. Sentry captures with full context
4. API returns 500 to user
5. You see it immediately in Sentry + Grafana

### Scenario: Someone is brute-forcing login

1. Rate limiter kicks in after 60 requests/minute
2. Returns 429 to attacker
3. `volo_rate_limit_hits_total` spikes
4. Grafana alert fires
5. You can block the IP if needed

---

## Files

| File | What it does |
|------|-------------|
| `internal/middleware/sentry.go` | Sentry SDK init, middleware, helpers |
| `internal/middleware/metrics.go` | Prometheus metric definitions + middleware |
| `internal/middleware/validate.go` | Input validation functions |
| `internal/middleware/validate_test.go` | 16 validation tests |
| `internal/middleware/logger.go` | Structured JSON request logging |
| `internal/middleware/recoverer.go` | Panic recovery → 500 JSON + Sentry |
| `cmd/server/main.go` | Wires everything: Sentry init, metrics middleware, `/metrics` endpoint |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VOLO_SENTRY_DSN` | (empty) | Sentry DSN. Empty = disabled. |
| `VOLO_ENVIRONMENT` | development | Sentry environment tag (development/production) |
| `VOLO_LOG_LEVEL` | info | Log verbosity (info/debug) |
