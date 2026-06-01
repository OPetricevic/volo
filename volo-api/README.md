# Volo API

Go backend for the Volo voice assistant. Handles intent parsing, command history, pattern learning, authentication, and rate limiting.

## Quick Start (Docker)

```bash
docker compose up -d
```

This starts:
- **Postgres** on :5432 (user: volo, pass: volo, db: volo)
- **Redis** on :6379
- **Volo API** on :8080

Migrations run automatically on first Postgres start.

## Local Development (without Docker for the API)

```bash
# Start only Postgres + Redis
docker compose up -d postgres redis

# Run the API directly
export VOLO_DATABASE_URL="postgres://volo:volo@localhost:5432/volo?sslmode=disable"
export VOLO_REDIS_ADDR="localhost:6379"
export VOLO_JWT_SECRET="dev-secret"
go run ./cmd/server
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VOLO_PORT` | 8080 | Server port |
| `VOLO_DATABASE_URL` | postgres://...localhost/volo | Postgres connection string |
| `VOLO_REDIS_ADDR` | localhost:6379 | Redis address |
| `VOLO_REDIS_PASSWORD` | (empty) | Redis password |
| `VOLO_JWT_SECRET` | dev-secret... | JWT signing secret |
| `VOLO_CORS_ORIGINS` | * | Comma-separated allowed origins |
| `VOLO_LOG_LEVEL` | info | Log level (info, debug) |
| `VOLO_GOOGLE_CLIENT_ID` | (empty) | Google OAuth client ID |
| `VOLO_GOOGLE_CLIENT_SECRET` | (empty) | Google OAuth client secret |

## Testing

```bash
# Unit tests (no Docker needed)
go test -short ./...

# All tests including integration (needs Docker running)
go test ./...
```

## Project Structure

```
cmd/server/main.go        Entry point, wiring
internal/
  config/                 Environment variable loading
  handler/                HTTP handlers (thin, validation)
  service/                Business logic
  repository/             Database queries (Postgres)
  intent/                 Rule-based command parser
  middleware/             Auth, CORS, rate limit, logging, request ID, recovery
  model/                  Domain structs, request/response types
  testutil/               Test container helpers
migrations/               SQL migration files
```

## API Overview

See [docs/API.md](./docs/API.md) for full reference.

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /api/v1/auth/device | No | Register device (anonymous) |
| POST | /api/v1/auth/google | No | Google SSO |
| POST | /api/v1/auth/register | No | Email + password signup |
| POST | /api/v1/auth/login | No | Email + password login |
| POST | /api/v1/command | JWT | Process voice command |
| GET | /api/v1/suggestions | JWT | Get autocomplete suggestions |
| GET | /api/v1/history | JWT | Get command history |
| DELETE | /api/v1/history | JWT | Clear history |
| GET | /api/v1/settings | JWT | Get user settings |
| PUT | /api/v1/settings | JWT | Update settings |
| POST | /api/v1/auth/logout | JWT | Revoke current session |
| POST | /api/v1/auth/logout-all | JWT | Revoke all sessions |
| GET | /health | No | Health check |
