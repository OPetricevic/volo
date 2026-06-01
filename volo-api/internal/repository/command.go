package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/volo/volo-api/internal/model"
)

type CommandRepository struct {
	pool *pgxpool.Pool
}

func NewCommandRepository(pool *pgxpool.Pool) *CommandRepository {
	return &CommandRepository{pool: pool}
}

func (r *CommandRepository) Create(ctx context.Context, cmd *model.Command) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO commands (id, user_id, transcript, parsed_action, parsed_target, parsed_query, confidence, executed_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		cmd.ID, cmd.UserID, cmd.Transcript, cmd.ParsedAction, cmd.ParsedTarget, cmd.ParsedQuery, cmd.Confidence, cmd.ExecutedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.Command.Create: %w", err)
	}
	return nil
}

func (r *CommandRepository) GetByUser(ctx context.Context, userID string, limit, offset int) ([]model.Command, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM commands WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("repository.Command.GetByUser (count): %w", err)
	}

	// Get paginated results
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, transcript, parsed_action, parsed_target, parsed_query, confidence, executed_at
		 FROM commands WHERE user_id = $1 ORDER BY executed_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("repository.Command.GetByUser (query): %w", err)
	}
	defer rows.Close()

	var commands []model.Command
	for rows.Next() {
		var cmd model.Command
		if err := rows.Scan(&cmd.ID, &cmd.UserID, &cmd.Transcript, &cmd.ParsedAction, &cmd.ParsedTarget, &cmd.ParsedQuery, &cmd.Confidence, &cmd.ExecutedAt); err != nil {
			return nil, 0, fmt.Errorf("repository.Command.GetByUser (scan): %w", err)
		}
		commands = append(commands, cmd)
	}

	return commands, total, nil
}

func (r *CommandRepository) DeleteByUser(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM commands WHERE user_id = $1`, userID,
	)
	if err != nil {
		return fmt.Errorf("repository.Command.DeleteByUser: %w", err)
	}
	return nil
}
