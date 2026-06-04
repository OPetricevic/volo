package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MacroAction represents a single step in a macro.
type MacroAction struct {
	Type string `json:"type"` // "navigate", "search", etc.
	URL  string `json:"url,omitempty"`
}

// Macro represents a user-defined custom voice command.
type Macro struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	TriggerPhrase string        `json:"trigger_phrase"`
	Name          string        `json:"name"`
	Actions       []MacroAction `json:"actions"`
	Enabled       bool          `json:"enabled"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type MacroRepository struct {
	pool *pgxpool.Pool
}

func NewMacroRepository(pool *pgxpool.Pool) *MacroRepository {
	return &MacroRepository{pool: pool}
}

const maxFreeMacros = 10

// List returns all active macros for a user.
func (r *MacroRepository) List(ctx context.Context, userID string) ([]Macro, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, trigger_phrase, name, actions, enabled, created_at, updated_at
		 FROM macros WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.Macro.List: %w", err)
	}
	defer rows.Close()

	var macros []Macro
	for rows.Next() {
		var m Macro
		var actionsJSON []byte
		if err := rows.Scan(&m.ID, &m.UserID, &m.TriggerPhrase, &m.Name, &actionsJSON, &m.Enabled, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("repository.Macro.List → Scan: %w", err)
		}
		if err := json.Unmarshal(actionsJSON, &m.Actions); err != nil {
			return nil, fmt.Errorf("repository.Macro.List → Unmarshal: %w", err)
		}
		macros = append(macros, m)
	}
	return macros, nil
}

// ListEnabled returns only enabled macros (used by extension for matching).
func (r *MacroRepository) ListEnabled(ctx context.Context, userID string) ([]Macro, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, trigger_phrase, name, actions, enabled, created_at, updated_at
		 FROM macros WHERE user_id = $1 AND deleted_at IS NULL AND enabled = true ORDER BY created_at ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.Macro.ListEnabled: %w", err)
	}
	defer rows.Close()

	var macros []Macro
	for rows.Next() {
		var m Macro
		var actionsJSON []byte
		if err := rows.Scan(&m.ID, &m.UserID, &m.TriggerPhrase, &m.Name, &actionsJSON, &m.Enabled, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("repository.Macro.ListEnabled → Scan: %w", err)
		}
		if err := json.Unmarshal(actionsJSON, &m.Actions); err != nil {
			return nil, fmt.Errorf("repository.Macro.ListEnabled → Unmarshal: %w", err)
		}
		macros = append(macros, m)
	}
	return macros, nil
}

// Count returns how many active macros a user has.
func (r *MacroRepository) Count(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM macros WHERE user_id = $1 AND deleted_at IS NULL`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository.Macro.Count: %w", err)
	}
	return count, nil
}

// Create adds a new macro. Returns error if user is at free limit.
func (r *MacroRepository) Create(ctx context.Context, userID, triggerPhrase, name string, actions []MacroAction) (*Macro, error) {
	count, err := r.Count(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.Macro.Create → Count: %w", err)
	}
	if count >= maxFreeMacros {
		return nil, fmt.Errorf("repository.Macro.Create: macro limit reached (max %d for free tier)", maxFreeMacros)
	}

	actionsJSON, err := json.Marshal(actions)
	if err != nil {
		return nil, fmt.Errorf("repository.Macro.Create → Marshal: %w", err)
	}

	id := uuid.New().String()
	now := time.Now()

	_, err = r.pool.Exec(ctx,
		`INSERT INTO macros (id, user_id, trigger_phrase, name, actions, enabled, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, true, $6, $6)`,
		id, userID, triggerPhrase, name, actionsJSON, now)
	if err != nil {
		return nil, fmt.Errorf("repository.Macro.Create → Insert: %w", err)
	}

	return &Macro{
		ID:            id,
		UserID:        userID,
		TriggerPhrase: triggerPhrase,
		Name:          name,
		Actions:       actions,
		Enabled:       true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Update modifies an existing macro.
func (r *MacroRepository) Update(ctx context.Context, userID, macroID, triggerPhrase, name string, actions []MacroAction, enabled bool) error {
	actionsJSON, err := json.Marshal(actions)
	if err != nil {
		return fmt.Errorf("repository.Macro.Update → Marshal: %w", err)
	}

	result, err := r.pool.Exec(ctx,
		`UPDATE macros SET trigger_phrase = $1, name = $2, actions = $3, enabled = $4, updated_at = now()
		 WHERE id = $5 AND user_id = $6 AND deleted_at IS NULL`,
		triggerPhrase, name, actionsJSON, enabled, macroID, userID)
	if err != nil {
		return fmt.Errorf("repository.Macro.Update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository.Macro.Update: macro not found")
	}
	return nil
}

// Delete soft-deletes a macro.
func (r *MacroRepository) Delete(ctx context.Context, userID, macroID string) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE macros SET deleted_at = now() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		macroID, userID)
	if err != nil {
		return fmt.Errorf("repository.Macro.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository.Macro.Delete: macro not found")
	}
	return nil
}
