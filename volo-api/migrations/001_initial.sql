-- 001_initial.sql
-- Initial schema for Volo API

-- Lookup tables (integer IDs, cached in Redis, never deleted)
CREATE TABLE action_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sites (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    base_url TEXT NOT NULL,
    search_url_template TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Core identity
CREATE TABLE users (
    id UUID PRIMARY KEY,
    display_name TEXT,
    email TEXT UNIQUE,
    avatar_url TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL AND deleted_at IS NULL;

-- Auth credentials (separate from users)
CREATE TABLE credentials (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_user_id TEXT,
    password_hash TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_credentials_user ON credentials(user_id);

-- Devices
CREATE TABLE devices (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id TEXT NOT NULL UNIQUE,
    device_name TEXT,
    platform TEXT,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_devices_user ON devices(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_device_id ON devices(device_id) WHERE deleted_at IS NULL;

-- Sessions (revocable JWTs)
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id),
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_token ON sessions(token_hash) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_user ON sessions(user_id) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_expires ON sessions(expires_at) WHERE revoked_at IS NULL;

-- Command history
CREATE TABLE commands (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transcript TEXT NOT NULL,
    parsed_action TEXT NOT NULL,
    parsed_target TEXT,
    parsed_query TEXT,
    confidence REAL NOT NULL DEFAULT 0,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_commands_user_time ON commands(user_id, executed_at DESC);

-- Pattern learning
CREATE TABLE patterns (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pattern_type TEXT NOT NULL,
    value TEXT NOT NULL,
    frequency INTEGER NOT NULL DEFAULT 1,
    last_used TIMESTAMPTZ NOT NULL DEFAULT now(),
    hour_weights JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, pattern_type, value)
);

CREATE INDEX idx_patterns_user_freq ON patterns(user_id, frequency DESC);

-- Audit logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    request_id TEXT NOT NULL,
    action TEXT NOT NULL,
    status TEXT NOT NULL,
    error_chain TEXT,
    metadata JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_logs(action, created_at DESC);

-- Seed lookup data
INSERT INTO action_types (name) VALUES
    ('search'),
    ('navigate'),
    ('open-and-search'),
    ('open-and-play'),
    ('browser-control');

INSERT INTO sites (name, base_url, search_url_template) VALUES
    ('youtube', 'https://www.youtube.com', 'https://www.youtube.com/results?search_query=%s'),
    ('google', 'https://www.google.com', 'https://www.google.com/search?q=%s'),
    ('github', 'https://www.github.com', 'https://github.com/search?q=%s'),
    ('spotify', 'https://open.spotify.com', 'https://open.spotify.com/search/%s'),
    ('reddit', 'https://www.reddit.com', 'https://www.reddit.com/search/?q=%s'),
    ('twitter', 'https://www.x.com', NULL),
    ('discord', 'https://discord.com', NULL),
    ('gmail', 'https://mail.google.com', NULL),
    ('netflix', 'https://www.netflix.com', NULL),
    ('twitch', 'https://www.twitch.tv', NULL),
    ('linkedin', 'https://www.linkedin.com', NULL),
    ('stackoverflow', 'https://stackoverflow.com', NULL);
