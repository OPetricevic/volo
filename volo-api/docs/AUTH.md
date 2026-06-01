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
| Token in response | Only returned once on login/register. Client stores securely. |
| CORS | Whitelist specific origins in production |
| Session fixation | New session created on every login |

## Adding New Auth Providers

To add GitHub OAuth (or any provider):

1. Add `VOLO_GITHUB_CLIENT_ID` + `VOLO_GITHUB_CLIENT_SECRET` to config
2. Create handler `POST /auth/github`
3. Exchange code for token, get user info
4. Upsert user + credential with `provider="github"`
5. Issue JWT

The `credentials` table supports unlimited providers per user.
