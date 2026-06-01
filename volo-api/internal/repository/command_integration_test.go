package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/testutil"
)

func TestCommandRepository_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	cmdRepo := NewCommandRepository(db.Pool)
	ctx := context.Background()

	// Create user first
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create commands
	target := "youtube"
	query := "lofi beats"
	for i := 0; i < 5; i++ {
		cmd := &model.Command{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			Transcript:   "open youtube lofi beats",
			ParsedAction: "open-and-search",
			ParsedTarget: &target,
			ParsedQuery:  &query,
			Confidence:   0.9,
			ExecutedAt:   now.Add(time.Duration(i) * time.Minute),
		}
		if err := cmdRepo.Create(ctx, cmd); err != nil {
			t.Fatalf("Create command %d failed: %v", i, err)
		}
	}

	// Get paginated
	commands, total, err := cmdRepo.GetByUser(ctx, user.ID, 3, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(commands) != 3 {
		t.Errorf("len(commands) = %d, want 3 (page size)", len(commands))
	}

	// Verify ordering (most recent first)
	if commands[0].ExecutedAt.Before(commands[1].ExecutedAt) {
		t.Error("expected commands ordered by executed_at DESC")
	}

	// Get page 2
	commands2, _, err := cmdRepo.GetByUser(ctx, user.ID, 3, 3)
	if err != nil {
		t.Fatalf("GetByUser page 2 failed: %v", err)
	}
	if len(commands2) != 2 {
		t.Errorf("page 2 len = %d, want 2", len(commands2))
	}
}

func TestCommandRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	cmdRepo := NewCommandRepository(db.Pool)
	ctx := context.Background()

	// Create user
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	userRepo.Create(ctx, user)

	// Create command
	cmd := &model.Command{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		Transcript:   "search test",
		ParsedAction: "search",
		Confidence:   0.9,
		ExecutedAt:   now,
	}
	cmdRepo.Create(ctx, cmd)

	// Delete
	if err := cmdRepo.DeleteByUser(ctx, user.ID); err != nil {
		t.Fatalf("DeleteByUser failed: %v", err)
	}

	// Verify empty
	commands, total, _ := cmdRepo.GetByUser(ctx, user.ID, 10, 0)
	if total != 0 {
		t.Errorf("total after delete = %d, want 0", total)
	}
	if len(commands) != 0 {
		t.Errorf("commands after delete = %d, want 0", len(commands))
	}
}
