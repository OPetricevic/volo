# API Reference

Base URL: `http://localhost:8080/api/v1`

All responses use the envelope format:
```json
// Success
{ "data": { ... } }

// Error
{ "error": { "code": "...", "message": "...", "internal_error": "..." } }
```

---

## Authentication

### POST /auth/device

Register a device (anonymous user). Called on first extension install.

**Request:**
```json
{
  "device_id": "ext-chrome-abc123",
  "device_name": "Chrome on Windows",
  "platform": "extension"
}
```

**Response (201):**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { "id": "uuid", "role": "user", "created_at": "..." },
    "device": { "id": "uuid", "device_id": "ext-chrome-abc123", ... }
  }
}
```

### POST /auth/register

Create account with email + password. Can upgrade an anonymous user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "device_id": "ext-chrome-abc123"
}
```

### POST /auth/login

Login with email + password.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "device_id": "ext-chrome-abc123"
}
```

### POST /auth/google

Google OAuth SSO. (Not yet implemented)

### POST /auth/logout 🔒

Revoke current session.

### POST /auth/logout-all 🔒

Revoke all sessions for the user.

---

## Commands

### POST /command 🔒

Process a voice command.

**Request:**
```json
{
  "transcript": "open youtube play lofi hip hop",
  "context": {
    "current_url": "https://google.com",
    "timestamp": "2026-06-01T08:35:00Z"
  }
}
```

**Response (200):**
```json
{
  "data": {
    "action": "open-and-play",
    "target": "youtube",
    "query": "lofi hip hop",
    "confidence": 0.92,
    "suggestions": [],
    "execute": {
      "url": "https://www.youtube.com/results?search_query=lofi+hip+hop",
      "auto_play": true
    }
  }
}
```

### GET /suggestions 🔒

Get autocomplete suggestions based on user history.

**Response (200):**
```json
{
  "data": {
    "suggestions": [
      { "text": "open youtube lofi", "type": "command", "score": 0.95 },
      { "text": "react hooks", "type": "search", "score": 0.8 }
    ]
  }
}
```

### GET /history 🔒

Get paginated command history.

**Query params:** `page` (default 1), `page_size` (default 20, max 100)

**Response (200):**
```json
{
  "data": {
    "commands": [
      {
        "id": "uuid",
        "transcript": "open youtube lofi",
        "parsed_action": "open-and-search",
        "parsed_target": "youtube",
        "parsed_query": "lofi",
        "confidence": 0.9,
        "executed_at": "2026-06-01T08:35:00Z"
      }
    ],
    "total": 42,
    "page": 1,
    "page_size": 20
  }
}
```

### DELETE /history 🔒

Clear all command history for the user.

---

## Settings

### GET /settings 🔒

Get user settings (synced across devices).

### PUT /settings 🔒

Update settings.

**Request:**
```json
{
  "mic_mode": "always",
  "wake_word": "hey volo",
  "language": "en-US"
}
```

---

## Health

### GET /health

**Response (200):**
```json
{ "data": { "status": "ok" } }
```

---

## Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| INVALID_REQUEST | 400 | Malformed request body |
| MISSING_DEVICE_ID | 400 | Device ID required |
| MISSING_FIELDS | 400 | Required fields missing |
| WEAK_PASSWORD | 400 | Password < 8 chars |
| EMPTY_TRANSCRIPT | 400 | Transcript cannot be empty |
| UNAUTHORIZED | 401 | Missing or invalid token |
| INVALID_CREDENTIALS | 401 | Wrong email/password |
| RATE_LIMITED | 429 | Too many requests (60/min) |
| INTERNAL_ERROR | 500 | Unexpected server error |
| REGISTRATION_FAILED | 500 | Failed to create account |
| COMMAND_PARSE_FAILED | 500 | Failed to process command |

---

🔒 = Requires `Authorization: Bearer <token>` header
