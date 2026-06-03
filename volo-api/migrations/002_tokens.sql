-- 002_tokens.sql
-- Generic tokens table for password resets, email verification, invites, etc.

CREATE TABLE tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,              -- 'password_reset', 'email_verify', 'invite'
    token_hash TEXT NOT NULL,        -- SHA-256 of the raw token
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,            -- NULL until consumed
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tokens_hash ON tokens(token_hash) WHERE used_at IS NULL;
CREATE INDEX idx_tokens_user_type ON tokens(user_id, type, created_at DESC);
