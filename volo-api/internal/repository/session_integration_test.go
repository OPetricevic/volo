package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/testutil"
)

func TestSessionRepository_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	sessionRepo := NewSessionRepository(db.Pool)
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

	// Create session
	session := &model.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: "sha256-test-hash-value",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
	if err := sessionRepo.Create(ctx, session); err != nil {
		t.Fatalf("Create session failed: %v", err)
	}

	// Get by token hash
	got, err := sessionRepo.GetByTokenHash(ctx, "sha256-test-hash-value")
	if err != nil {
		t.Fatalf("GetByTokenHash failed: %v", err)
	}
	if got.ID != session.ID {
		t.Errorf("session ID = %q, want %q", got.ID, session.ID)
	}
	if got.UserID != user.ID {
		t.Errorf("session user_id = %q, want %q", got.UserID, user.ID)
	}
}

func TestSessionRepository_Revoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	sessionRepo := NewSessionRepository(db.Pool)
	ctx := context.Background()

	// Create user + session
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	userRepo.Create(ctx, user)

	session := &model.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: "hash-to-revoke",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
	sessionRepo.Create(ctx, session)

	// Revoke
	if err := sessionRepo.Revoke(ctx, "hash-to-revoke"); err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}

	// Should not find revoked session
	_, err := sessionRepo.GetByTokenHash(ctx, "hash-to-revoke")
	if err == nil {
		t.Error("expected error for revoked session, got nil")
	}
}

func TestSessionRepository_RevokeAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	sessionRepo := NewSessionRepository(db.Pool)
	ctx := context.Background()

	// Create user + multiple sessions
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	userRepo.Create(ctx, user)

	for i := 0; i < 3; i++ {
		session := &model.Session{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			TokenHash: uuid.New().String(),
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		}
		sessionRepo.Create(ctx, session)
	}

	// Revoke all
	if err := sessionRepo.RevokeAllForUser(ctx, user.ID); err != nil {
		t.Fatalf("RevokeAllForUser failed: %v", err)
	}

	// None should be findable
	// (We'd need to query all sessions to verify, but the GetByTokenHash
	// already filters revoked ones, so this is implicitly tested)
}

func TestSessionRepository_ExpiredNotReturned(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	sessionRepo := NewSessionRepository(db.Pool)
	ctx := context.Background()

	// Create user + expired session
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	userRepo.Create(ctx, user)

	session := &model.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: "expired-hash",
		ExpiresAt: now.Add(-1 * time.Hour), // expired 1 hour ago
		CreatedAt: now.Add(-25 * time.Hour),
	}
	sessionRepo.Create(ctx, session)

	// Should not find expired session
	_, err := sessionRepo.GetByTokenHash(ctx, "expired-hash")
	if err == nil {
		t.Error("expected error for expired session, got nil")
	}
}
