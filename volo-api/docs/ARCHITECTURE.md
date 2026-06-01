# API Architecture

## Layered Design

```
HTTP Request
    │
    ▼
┌─────────────────────────────────────────────┐
│  Middleware Stack                             │
│  RealIP → RequestID → Logger → Recoverer    │
│  → CORS → Auth → RateLimit                  │
└──────────────────────┬──────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────┐
│  Handler Layer                               │
│  • Decode request body                       │
│  • Validate input                            │
│  • Call service                              │
│  • Format response                           │
│  • NO business logic here                    │
└──────────────────────┬──────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────┐
│  Service Layer                               │
│  • Business logic                            │
│  • UUID generation                           │
│  • Password hashing                          │
│  • JWT creation/validation                   │
│  • Orchestrates repositories                 │
└──────────────────────┬──────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────┐
│  Repository Layer                            │
│  • SQL queries only                          │
│  • One method = one query                    │
│  • Returns domain models                     │
│  • Error wrapping with location              │
└──────────────────────┬──────────────────────┘
                       │
              ┌────────┴────────┐
              ▼                 ▼
        ┌──────────┐     ┌──────────┐
        │ Postgres │     │  Redis   │
        └──────────┘     └──────────┘
```

## Error Handling

Every error is wrapped with its location for debugging:

```go
// Repository
return fmt.Errorf("repository.User.GetByID: %w", err)

// Service
return fmt.Errorf("service.Auth.RegisterDevice → GetByID: %w", err)

// Handler
respondError(w, 500, "REGISTRATION_FAILED",
    "Failed to register device. Please try again.",
    err.Error())  // full chain goes in internal_error
```

Response format (both friendly + debug info):
```json
{
  "error": {
    "code": "REGISTRATION_FAILED",
    "message": "Failed to register device. Please try again.",
    "internal_error": "service.Auth.RegisterDevice → repository.User.GetByID: no rows"
  }
}
```

## Request Flow Example

`POST /api/v1/command` with `{"transcript": "open youtube lofi"}`

1. **Middleware:** Validates JWT, extracts user_id, checks rate limit
2. **Handler:** Decodes body, validates transcript not empty
3. **Service:** Calls intent parser, stores command in history, builds response
4. **Repository:** Inserts into `commands` table
5. **Response:** Returns parsed intent with execute URL

## Graceful Shutdown

```
SIGTERM received
    → Stop accepting new connections
    → Wait up to 10s for in-flight requests
    → Close database pool
    → Close Redis client
    → Exit
```

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| UUIDs generated in Go, not Postgres | Testable, no pgcrypto dependency |
| Passwords in `credentials` table, not `users` | Security separation, supports multiple auth methods |
| Sessions stored in DB + cached in Redis | Revocable (DB) + fast validation (Redis) |
| Rate limiting via Redis sorted sets | Sliding window, accurate, survives restarts |
| Lookup tables with integer IDs | Fast joins, cached in Redis, never deleted |
| Soft deletes (`deleted_at`) | Audit trail, recoverable, no FK violations |
