package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/volo/volo-api/internal/model"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	metadata, _ := json.Marshal(log.Metadata)

	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_logs (id, user_id, request_id, action, status, error_chain, metadata, ip_address, user_agent, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		log.ID, log.UserID, log.RequestID, log.Action, log.Status, log.ErrorChain, metadata, log.IPAddress, log.UserAgent, log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.Audit.Create: %w", err)
	}
	return nil
}
