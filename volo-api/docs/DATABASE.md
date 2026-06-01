# Database Schema

## Overview

PostgreSQL with two categories of tables:

- **Lookup tables** — Integer serial IDs, never deleted, cached in Redis
- **Core tables** — UUID primary keys (generated in Go), soft deletes, full timestamps

## Lookup Tables

### action_types

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL PK | Auto-increment |
| name | TEXT UNIQUE | search, navigate, open-and-search, open-and-play, browser-control |
| created_at | TIMESTAMPTZ | Default now() |

### sites

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL PK | Auto-increment |
| name | TEXT UNIQUE | youtube, google, github, etc. |
| base_url | TEXT | https://www.youtube.com |
| search_url_template | TEXT | https://youtube.com/results?search_query=%s |
| created_at | TIMESTAMPTZ | Default now() |

## Core Tables

### users

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | Generated in Go |
| display_name | TEXT | Nullable |
| email | TEXT UNIQUE | Nullable until account upgrade |
| avatar_url | TEXT | Nullable |
| role | TEXT | Default 'user' |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |
| deleted_at | TIMESTAMPTZ | Soft delete |

### credentials

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK → users | ON DELETE CASCADE |
| provider | TEXT | 'password', 'google' |
| provider_user_id | TEXT | Google sub ID (null for password) |
| password_hash | TEXT | bcrypt (null for OAuth) |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

Unique constraint: `(provider, provider_user_id)`

### devices

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK → users | |
| device_id | TEXT UNIQUE | Client-generated identifier |
| device_name | TEXT | "Chrome on Windows" |
| platform | TEXT | 'extension', 'desktop' |
| last_seen_at | TIMESTAMPTZ | Updated on each request |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |
| deleted_at | TIMESTAMPTZ | Soft delete |

### sessions

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK → users | |
| device_id | UUID FK → devices | Nullable |
| token_hash | TEXT UNIQUE | SHA256 of JWT |
| expires_at | TIMESTAMPTZ | 30 days from creation |
| revoked_at | TIMESTAMPTZ | Set on logout |
| created_at | TIMESTAMPTZ | |

### commands

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK → users | |
| transcript | TEXT | Raw voice input |
| parsed_action | TEXT | search, navigate, etc. |
| parsed_target | TEXT | Nullable (youtube, github...) |
| parsed_query | TEXT | Nullable (search terms) |
| confidence | REAL | 0.0 – 1.0 |
| executed_at | TIMESTAMPTZ | |

### patterns

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK → users | |
| pattern_type | TEXT | 'search', 'site', 'command' |
| value | TEXT | The pattern value |
| frequency | INTEGER | Times used |
| last_used | TIMESTAMPTZ | |
| hour_weights | JSONB | {"8": 5, "9": 3} |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

Unique constraint: `(user_id, pattern_type, value)`

### audit_logs

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK → users | Nullable (failed auth) |
| request_id | TEXT | From X-Request-ID header |
| action | TEXT | 'command.process', 'auth.login', etc. |
| status | TEXT | 'success', 'error' |
| error_chain | TEXT | Full wrapped error path |
| metadata | JSONB | Request context |
| ip_address | INET | |
| user_agent | TEXT | |
| created_at | TIMESTAMPTZ | |

## Migrations

Located in `migrations/`. Run sequentially on startup.

| File | Description |
|------|-------------|
| 001_initial.sql | Full schema + seed data |

To add a new migration: create `002_description.sql` with ALTER/CREATE statements.

## Indexes

| Table | Index | Purpose |
|-------|-------|---------|
| users | email (WHERE NOT NULL, NOT deleted) | Login lookup |
| credentials | user_id | Find auth methods for user |
| devices | user_id (WHERE NOT deleted) | List user's devices |
| devices | device_id (WHERE NOT deleted) | Device registration check |
| sessions | token_hash (WHERE NOT revoked) | Token validation |
| sessions | user_id (WHERE NOT revoked) | Revoke all |
| commands | user_id, executed_at DESC | History pagination |
| patterns | user_id, frequency DESC | Top patterns |
| audit_logs | user_id, created_at DESC | User audit trail |
| audit_logs | action, created_at DESC | Action-based queries |
