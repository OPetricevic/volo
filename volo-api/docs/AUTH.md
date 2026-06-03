# Authentication

## Overview

Volo uses a progressive auth model:

1. **Anonymous** — Zero friction on install. Device gets a JWT immediately.
2. **Account upgrade** — User can optionally add email+password or Google SSO.
3. **Multi-device** — Account links multiple devices. History syncs.

## Flows

### 1. Fresh Install (Anonymous)

```
Extension installs
    → Generates random device_id locally
    → POST /auth/device { device_id, device_name, platform }
    → Backend creates: user (no email) + device + session
    → Returns JWT
    → Extension stores JWT in chrome.storage.local
    → Ready to use immediately
```

No signup wall. No email. No friction.

### 2. Google SSO (Upgrade or Login)

```
User clicks "Sign in with Google" in settings
    → Extension opens Google OAuth consent screen
    → User grants access → gets authorization code
    → POST /auth/google { code, device_id }
    → Backend exchanges code for Google token
    → Gets email + Google sub ID
    → If existing user with that Google sub → link device, return JWT
    → If new → create user with email, add credential(google), return JWT
    → If anonymous upgrading → update user email, add credential(google)
```

### 3. Email + Password

```
POST /auth/register { email, password, device_id }
    → If user is authenticated (upgrading anonymous):
        → Update user email
        → Create credential(password) with bcrypt hash
    → If new user:
        → Create user + credential + device
    → Return JWT

POST /auth/login { email, password, device_id }
    → Find user by email
    → Find credential(password) for user
    → Compare bcrypt hash
    → Link device if new
    → Create session, return JWT
```

### 4. Logout

```
POST /auth/logout
    → Revoke current session (set revoked_at)
    → Delete Redis cache entry
    → Token immediately invalid

POST /auth/logout-all
    → Revoke ALL sessions for user
    → All devices logged out
```

### 5. Password Reset

```
POST /auth/forgot-password { email }
    → Find user by email (if not found, return 200 anyway — no enumeration)
    → Rate limit: max 3 resets per hour per user
    → Generate 32-byte random token
    → Store SHA-256 hash in tokens table (type = "password_reset", expires in 1 hour)
    → Send email via Resend with reset link
    → Return 200 { "message": "If an account with that email exists, a reset link has been sent." }

POST /auth/reset-password { token, password }
    → Hash the raw token → look up in tokens table (unused + not expired)
    → Validate new password (min 8 chars, max 128)
    → Update credentials table (bcrypt hash)
    → Mark token as used
    → Invalidate all other reset tokens for this user
    → Revoke ALL sessions (force re-login everywhere)
    → Return 200 { "message": "Password reset successful." }
```

**Security:**
- Tokens are SHA-256 hashed in DB (raw token only exists in the email link)
- 1-hour expiry, single-use
- Rate limited at the application level (3 requests per email per hour)
- Response always returns 200 on forgot-password (prevents email enumeration)
- All sessions revoked after reset (attacker can't stay logged in)

**Email delivery:**
- Uses [Resend](https://resend.com) (3,000 emails/month free tier)
- Configured via `VOLO_RESEND_API_KEY` env var
- If not configured, email step is skipped (token still generated, logged as warning)
- Reset link points to `VOLO_RESET_PASSWORD_URL` (default: `http://localhost:5173/reset`)

**Database:**
- Generic `tokens` table supports multiple token types (password_reset, email_verify, invite)
- Indexed on `token_hash` for fast lookup
- Periodic cleanup: unused expired tokens can be purged via cron

## JWT Structure

```json
{
  "sub": "user-uuid",
  "iat": 1717228500,
  "exp": 1719820500
}
```

- Signed with HMAC-SHA256
- 30-day expiry
- Token hash (SHA256) stored in `sessions` table for revocation

## Token Validation Flow

```
Request arrives with Authorization: Bearer <token>
    │
    ▼
Parse JWT, verify signature + expiry
    │
    ▼
Compute SHA256(token) → token_hash
    │
    ▼
Check Redis: "session:<hash>" → user_id?
    ├── HIT: return user_id (fast path)
    │
    └── MISS: query sessions table
              WHERE token_hash = ? AND revoked_at IS NULL AND expires_at > now()
              │
              ├── FOUND: cache in Redis (TTL = remaining expiry), return user_id
              └── NOT FOUND: return 401 Unauthorized
```

## Security Considerations

| Concern | Mitigation |
|---------|-----------|
| Password storage | bcrypt with default cost (10 rounds) |
| Passwords in users table | Never. Separate `credentials` table. |
| Token theft | Sessions are revocable. Logout invalidates immediately. |
| Brute force | Rate limiting (60 req/min per user via Redis) |
| Password reset abuse | 3 resets/hour per user, hashed tokens, 1-hour expiry |
| Email enumeration | Forgot-password always returns 200 regardless |
| Token in response | Only returned once on login/register. Client stores securely. |
| CORS | Whitelist specific origins in production |
| Session fixation | New session created on every login |
| Post-compromise | Password reset revokes ALL sessions |

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `VOLO_JWT_SECRET` | Yes | HMAC signing secret (generate with `openssl rand -hex 32`) |
| `VOLO_GOOGLE_CLIENT_ID` | No | Google OAuth client ID |
| `VOLO_GOOGLE_CLIENT_SECRET` | No | Google OAuth client secret |
| `VOLO_RESEND_API_KEY` | No | Resend API key for password reset emails |
| `VOLO_RESET_PASSWORD_URL` | No | Frontend URL for reset page (default: `http://localhost:5173/reset`) |

## Adding New Auth Providers

To add GitHub OAuth (or any provider):

1. Add `VOLO_GITHUB_CLIENT_ID` + `VOLO_GITHUB_CLIENT_SECRET` to config
2. Create handler `POST /auth/github`
3. Exchange code for token, get user info
4. Upsert user + credential with `provider="github"`
5. Issue JWT

The `credentials` table supports unlimited providers per user.
