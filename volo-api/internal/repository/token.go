package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/volo/volo-api/internal/model"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

// Create inserts a new token.
func (r *TokenRepository) Create(ctx context.Context, token *model.Token) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tokens (id, user_id, type, token_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		token.ID, token.UserID, token.Type, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.Token.Create: %w", err)
	}
	return nil
}

// GetByHash retrieves a valid (unused, not expired) token by its hash.
func (r *TokenRepository) GetByHash(ctx context.Context, tokenHash string) (*model.Token, error) {
	var token model.Token
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, type, token_hash, expires_at, used_at, created_at
		 FROM tokens
		 WHERE token_hash = $1 AND used_at IS NULL AND expires_at > now()`,
		tokenHash,
	).Scan(&token.ID, &token.UserID, &token.Type, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.Token.GetByHash: %w", err)
	}
	return &token, nil
}

// MarkUsed marks a token as consumed.
func (r *TokenRepository) MarkUsed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tokens SET used_at = now() WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("repository.Token.MarkUsed: %w", err)
	}
	return nil
}

// CountRecentByUserAndType counts tokens created in the given window (for rate limiting).
func (r *TokenRepository) CountRecentByUserAndType(ctx context.Context, userID, tokenType string, since time.Time) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tokens WHERE user_id = $1 AND type = $2 AND created_at > $3`,
		userID, tokenType, since,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository.Token.CountRecentByUserAndType: %w", err)
	}
	return count, nil
}

// InvalidateAllForUser marks all unused tokens of a type as used (e.g., after successful reset).
func (r *TokenRepository) InvalidateAllForUser(ctx context.Context, userID, tokenType string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tokens SET used_at = now() WHERE user_id = $1 AND type = $2 AND used_at IS NULL`,
		userID, tokenType,
	)
	if err != nil {
		return fmt.Errorf("repository.Token.InvalidateAllForUser: %w", err)
	}
	return nil
}
