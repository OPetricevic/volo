# API Testing

## Overview

Two test tiers:

1. **Unit tests** — Fast, no external dependencies, run everywhere
2. **Integration tests** — Use testcontainers (ephemeral Docker), test real Postgres queries

## Running Tests

```bash
# Unit tests only (fast, no Docker)
go test -short ./...

# All tests (needs Docker running)
go test ./...

# Verbose output
go test -v ./...

# Single package
go test -v ./internal/intent/...
```

## Unit Tests (46 tests)

### Intent Parser (28 tests)

| Category | Tests | What's verified |
|----------|-------|----------------|
| Search | 5 | search, search for, look up, find, fallback |
| Navigate | 7 | open site, go to, goto, fuzzy matching, raw URLs |
| Open + search | 3 | youtube, spotify, github |
| Open + play | 2 | youtube play, spotify play |
| Browser control | 10 | back, forward, close tab, new tab, scroll |
| Case insensitivity | 2 | uppercase, mixed case |
| URL building | 4 | google, youtube, spotify, raw |

### Middleware (15 tests)

| Middleware | Tests | What's verified |
|-----------|-------|----------------|
| Request ID | 3 | Sets header, unique per request, missing context fallback |
| CORS | 5 | Wildcard, specific origin, reject unknown, preflight 204, headers |
| Auth | 5 | Valid token, missing header, invalid format, invalid token, missing context |
| Recoverer | 2 | Catches panic → 500 JSON, passes through normally |

### Handler (4 tests)

- Health endpoint returns 200 + `{"status":"ok"}`
- Error response format includes code + message + internal_error
- JSON decoding works
- Invalid JSON returns error

### Config (4 tests)

- Defaults load correctly
- Environment variables override defaults
- Wildcard CORS origin
- Multiple comma-separated origins

### Model (4 tests)

- Success response JSON format
- Error response JSON format
- Password hash never appears in JSON (json:"-" tag)
- Token hash never appears in JSON

## Integration Tests (11 tests, need Docker)

Uses `testcontainers-go` — spins up ephemeral Postgres containers per test. Zero state leaks.

### User Repository (3 tests)

- Create user + get by ID + get by email
- Create credential + get by user+provider
- Create device + get by device_id + update last_seen + soft delete

### Command Repository (2 tests)

- Create 5 commands + paginated retrieval (page 1: 3, page 2: 2) + ordering
- Delete all commands for user + verify empty

### Session Repository (4 tests)

- Create session + get by token hash
- Revoke session → not findable
- Revoke all for user
- Expired sessions not returned

### Audit Repository (2 tests)

- Create success + error audit logs with full metadata
- Create audit log without user (failed auth attempt)

## Test Infrastructure

```go
// testutil.SetupPostgres(t) does:
// 1. Starts postgres:16-alpine container
// 2. Runs 001_initial.sql migration
// 3. Returns connection pool
// 4. Registers cleanup (close pool + terminate container)

db := testutil.SetupPostgres(t)
repo := NewUserRepository(db.Pool)
// ... test queries against real Postgres
// Container destroyed automatically when test ends
```

## What's NOT Tested

| Area | Why | How to test |
|------|-----|-------------|
| Service layer (auth flow) | Needs mocked repos or full stack | Add service_test.go with mock interfaces |
| Rate limiting | Needs Redis container | Add Redis testcontainer |
| Full HTTP integration | Needs running server | httptest.Server + real DB |
| Google OAuth | Needs Google API | Mock HTTP client |
