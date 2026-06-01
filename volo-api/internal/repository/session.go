package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/volo/volo-api/internal/model"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sessions (id, user_id, device_id, token_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		session.ID, session.UserID, session.DeviceID, session.TokenHash, session.ExpiresAt, session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.Session.Create: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.Session, error) {
	var session model.Session
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, device_id, token_hash, expires_at, revoked_at, created_at
		 FROM sessions WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now()`,
		tokenHash,
	).Scan(&session.ID, &session.UserID, &session.DeviceID, &session.TokenHash, &session.ExpiresAt, &session.RevokedAt, &session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.Session.GetByTokenHash: %w", err)
	}
	return &session, nil
}

func (r *SessionRepository) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sessions SET revoked_at = now() WHERE token_hash = $1`, tokenHash,
	)
	if err != nil {
		return fmt.Errorf("repository.Session.Revoke: %w", err)
	}
	return nil
}

func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL`, userID,
	)
	if err != nil {
		return fmt.Errorf("repository.Session.RevokeAllForUser: %w", err)
	}
	return nil
}
