-- 002_macros.sql
-- Custom voice command macros (user-defined shortcuts)

CREATE TABLE macros (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trigger_phrase TEXT NOT NULL,
    name TEXT NOT NULL,
    actions JSONB NOT NULL DEFAULT '[]',
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, trigger_phrase)
);

CREATE INDEX idx_macros_user ON macros(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_macros_trigger ON macros(user_id, trigger_phrase) WHERE deleted_at IS NULL AND enabled = true;

-- actions JSONB format:
-- [
--   { "type": "navigate", "url": "https://gmail.com" },
--   { "type": "navigate", "url": "https://youtube.com" },
--   { "type": "navigate", "url": "https://weather.com/vienna" }
-- ]
--
-- Free tier: max 3 macros per user (enforced in application layer)
