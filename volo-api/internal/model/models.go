package model

import "time"

// User represents a user identity.
type User struct {
	ID          string     `json:"id"`
	DisplayName *string    `json:"display_name,omitempty"`
	Email       *string    `json:"email,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

// Credential represents an auth method for a user.
type Credential struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Provider       string    `json:"provider"` // "password", "google"
	ProviderUserID *string   `json:"provider_user_id,omitempty"`
	PasswordHash   *string   `json:"-"` // never serialized
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Device represents a linked device.
type Device struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	DeviceID   string     `json:"device_id"`
	DeviceName *string    `json:"device_name,omitempty"`
	Platform   *string    `json:"platform,omitempty"`
	LastSeenAt time.Time  `json:"last_seen_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-"`
}

// Session represents an active JWT session.
type Session struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	DeviceID  *string    `json:"device_id,omitempty"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}

// Command represents a processed voice command.
type Command struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Transcript   string    `json:"transcript"`
	ParsedAction string    `json:"parsed_action"`
	ParsedTarget *string   `json:"parsed_target,omitempty"`
	ParsedQuery  *string   `json:"parsed_query,omitempty"`
	Confidence   float64   `json:"confidence"`
	ExecutedAt   time.Time `json:"executed_at"`
}

// Pattern represents a learned user behavior pattern.
type Pattern struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	PatternType string         `json:"pattern_type"` // "search", "site", "command"
	Value       string         `json:"value"`
	Frequency   int            `json:"frequency"`
	LastUsed    time.Time      `json:"last_used"`
	HourWeights map[string]int `json:"hour_weights"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Token represents a generic token (password reset, email verify, invite).
type Token struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Type      string     `json:"type"` // "password_reset", "email_verify", "invite"
	TokenHash string     `json:"-"`    // never serialized
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// AuditLog represents a tracked action.
type AuditLog struct {
	ID         string                 `json:"id"`
	UserID     *string                `json:"user_id,omitempty"`
	RequestID  string                 `json:"request_id"`
	Action     string                 `json:"action"`
	Status     string                 `json:"status"` // "success", "error"
	ErrorChain *string                `json:"error_chain,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	IPAddress  *string                `json:"ip_address,omitempty"`
	UserAgent  *string                `json:"user_agent,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// Lookup tables

type ActionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Site struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	BaseURL           string  `json:"base_url"`
	SearchURLTemplate *string `json:"search_url_template,omitempty"`
}
