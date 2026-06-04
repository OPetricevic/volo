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

Google OAuth SSO.

**Request:**
```json
{
  "code": "google-auth-code-from-popup",
  "device_id": "ext-chrome-abc123"
}
```

### POST /auth/forgot-password

Request a password reset email. Always returns 200 (doesn't reveal if email exists).

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response (200):**
```json
{
  "data": { "message": "If an account with that email exists, a reset link has been sent." }
}
```

### POST /auth/reset-password

Reset password using the token from the email link.

**Request:**
```json
{
  "token": "base64-encoded-reset-token",
  "password": "newpassword123"
}
```

**Response (200):**
```json
{
  "data": { "message": "Password reset successful. Please log in with your new password." }
}
```

**Errors:** `MISSING_TOKEN`, `WEAK_PASSWORD`, `PASSWORD_TOO_LONG`, `RESET_FAILED`

### POST /auth/verify-email

Verify email address using the token from the verification email.

**Request:**
```json
{
  "token": "base64-encoded-verification-token"
}
```

**Response (200):**
```json
{
  "data": { "message": "Email verified successfully." }
}
```

**Errors:** `MISSING_TOKEN`, `VERIFICATION_FAILED`

### POST /auth/change-password 🔒

Change password while logged in. Requires current password.

**Request:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response (200):**
```json
{
  "data": { "message": "Password changed successfully." }
}
```

**Errors:** `MISSING_CURRENT_PASSWORD`, `WEAK_PASSWORD`, `PASSWORD_TOO_LONG`, `SAME_PASSWORD`, `CHANGE_PASSWORD_FAILED`

### POST /auth/resend-verification 🔒

Resend the email verification link. Rate limited to 3 per hour.

**Response (200):**
```json
{
  "data": { "message": "Verification email sent." }
}
```

**Errors:** `NO_EMAIL`, `ALREADY_VERIFIED`, `VERIFICATION_SEND_FAILED`

### DELETE /auth/account 🔒

Delete account (soft delete). Requires password confirmation.

**Request:**
```json
{
  "password": "currentpassword123"
}
```

**Response (200):**
```json
{
  "data": { "message": "Account deleted. We're sorry to see you go." }
}
```

**Errors:** `MISSING_PASSWORD`, `DELETE_FAILED`

### POST /auth/logout 🔒

Revoke current session.

### POST /auth/logout-all 🔒

Revoke all sessions for the user.

### DELETE /auth/device/{deviceID} 🔒

Unlink a device from the account. Only the device owner can unlink.

**Errors:** `MISSING_DEVICE_ID`, `DEVICE_NOT_FOUND`, `UNLINK_FAILED`

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

### GET /history/context 🔒

Get a formatted text summary of recent commands (used internally for LLM context).

---

## Chat (Ollama AI)

### POST /chat 🔒

Send a message to Volo AI (powered by local Ollama + Gemma 2B).

**Request:**
```json
{
  "message": "What did I search yesterday?",
  "conversation_id": "uuid (optional)"
}
```

**Response (200):**
```json
{
  "data": {
    "reply": "Yesterday you searched for:\n• react hooks tutorial (9:14am)\n• golang concurrency (2:30pm)",
    "sources": [
      { "id": "cmd-uuid", "transcript": "search react hooks tutorial", "executed_at": "2026-06-01T09:14:00Z" }
    ]
  }
}
```

**Errors:** `EMPTY_MESSAGE`, `MESSAGE_TOO_LONG`, `CHAT_FAILED` (Ollama not running)

### GET /chat/status 🔒

Check if Ollama is running and available.

**Response (200):**
```json
{
  "data": { "ollama_running": true, "model": "gemma2:2b" }
}
```

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

## Macros (Custom Voice Commands)

### GET /macros 🔒

List all macros for the authenticated user.

**Response (200):**
```json
{
  "data": {
    "macros": [
      {
        "id": "uuid",
        "trigger_phrase": "morning routine",
        "name": "Morning Routine",
        "actions": [
          { "type": "navigate", "url": "https://gmail.com" },
          { "type": "navigate", "url": "https://youtube.com" },
          { "type": "navigate", "url": "https://weather.com" }
        ],
        "enabled": true,
        "created_at": "2026-06-04T12:00:00Z",
        "updated_at": "2026-06-04T12:00:00Z"
      }
    ],
    "limit": 10,
    "count": 1
  }
}
```

### GET /macros/enabled 🔒

List only enabled macros (used by extension for syncing).

**Response (200):**
```json
{
  "data": {
    "macros": [ ... ]
  }
}
```

### POST /macros 🔒

Create a new macro.

**Request:**
```json
{
  "trigger_phrase": "morning routine",
  "name": "Morning Routine",
  "actions": [
    { "type": "navigate", "url": "https://gmail.com" },
    { "type": "navigate", "url": "https://youtube.com" }
  ]
}
```

**Response (201):**
```json
{
  "data": {
    "id": "uuid",
    "trigger_phrase": "morning routine",
    "name": "Morning Routine",
    "actions": [ ... ],
    "enabled": true,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

**Errors:** `MISSING_TRIGGER`, `MISSING_NAME`, `MISSING_ACTIONS`, `TOO_MANY_ACTIONS` (max 3), `MACRO_LIMIT` (max 10 per user)

### PUT /macros/{id} 🔒

Update an existing macro.

**Request:**
```json
{
  "trigger_phrase": "morning routine",
  "name": "Morning Routine (updated)",
  "actions": [ ... ],
  "enabled": false
}
```

**Errors:** `MISSING_FIELDS`, `MACRO_UPDATE_FAILED`

### DELETE /macros/{id} 🔒

Soft-delete a macro.

**Response (200):**
```json
{
  "data": { "status": "deleted" }
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
| MISSING_EMAIL | 400 | Email required |
| INVALID_EMAIL | 400 | Invalid email format |
| WEAK_PASSWORD | 400 | Password < 8 chars |
| PASSWORD_TOO_LONG | 400 | Password > 128 chars |
| SAME_PASSWORD | 400 | New password same as current |
| MISSING_TOKEN | 400 | Token required |
| MISSING_CURRENT_PASSWORD | 400 | Current password required |
| MISSING_PASSWORD | 400 | Password required for confirmation |
| EMPTY_TRANSCRIPT | 400 | Transcript cannot be empty |
| TRANSCRIPT_TOO_LONG | 400 | Transcript > 500 chars |
| EMPTY_MESSAGE | 400 | Chat message cannot be empty |
| MESSAGE_TOO_LONG | 400 | Chat message > 2000 chars |
| NO_EMAIL | 400 | No email on account |
| ALREADY_VERIFIED | 400 | Email already verified |
| DEVICE_NAME_TOO_LONG | 400 | Device name > 100 chars |
| UNAUTHORIZED | 401 | Missing or invalid token |
| INVALID_CREDENTIALS | 401 | Wrong email/password |
| DEVICE_NOT_FOUND | 404 | Device not found or not owned |
| RATE_LIMITED | 429 | Too many requests (60/min) |
| GOOGLE_NOT_CONFIGURED | 503 | Server missing Google OAuth config |
| REGISTRATION_FAILED | 500 | Failed to create account |
| COMMAND_PARSE_FAILED | 500 | Failed to process command |
| RESET_FAILED | 400 | Invalid/expired reset token |
| VERIFICATION_FAILED | 400 | Invalid/expired verification token |
| CHANGE_PASSWORD_FAILED | 400 | Current password incorrect |
| DELETE_FAILED | 400 | Password incorrect for deletion |
| CHAT_FAILED | 500 | Ollama not running |
| MISSING_TRIGGER | 400 | Macro trigger phrase required |
| MISSING_NAME | 400 | Macro name required |
| MISSING_ACTIONS | 400 | At least one action required |
| TOO_MANY_ACTIONS | 400 | Max 3 actions per macro |
| MACRO_LIMIT | 403 | Max 10 macros per user |
| MACRO_LIST_FAILED | 500 | Failed to load macros |
| MACRO_CREATE_FAILED | 500 | Failed to create macro |
| MACRO_UPDATE_FAILED | 500 | Failed to update macro |
| MACRO_DELETE_FAILED | 500 | Failed to delete macro |

---

🔒 = Requires `Authorization: Bearer <token>` header
