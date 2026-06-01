package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserSettings represents synced user preferences.
type UserSettings struct {
	MicMode  string `json:"mic_mode"`
	WakeWord string `json:"wake_word"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

type SettingsRepository struct {
	pool *pgxpool.Pool
}

func NewSettingsRepository(pool *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{pool: pool}
}

// Get retrieves user settings from the users.settings JSONB column.
// Returns defaults if no settings are stored.
func (r *SettingsRepository) Get(ctx context.Context, userID string) (*UserSettings, error) {
	var settingsJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(settings, '{}')::text FROM users WHERE id = $1 AND deleted_at IS NULL`, userID,
	).Scan(&settingsJSON)
	if err != nil {
		return nil, fmt.Errorf("repository.Settings.Get: %w", err)
	}

	settings := &UserSettings{
		MicMode:  "always",
		WakeWord: "hey volo",
		Language: "en-US",
		Theme:    "dark",
	}

	if len(settingsJSON) > 2 { // more than "{}"
		json.Unmarshal(settingsJSON, settings)
	}

	return settings, nil
}

// Update stores user settings in the users.settings JSONB column.
func (r *SettingsRepository) Update(ctx context.Context, userID string, settings *UserSettings) error {
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("repository.Settings.Update → marshal: %w", err)
	}

	_, err = r.pool.Exec(ctx,
		`UPDATE users SET settings = $1::jsonb, updated_at = now() WHERE id = $2`,
		string(settingsJSON), userID,
	)
	if err != nil {
		return fmt.Errorf("repository.Settings.Update: %w", err)
	}

	return nil
}
