package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/testutil"
)

func TestAuditRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	userRepo := NewUserRepository(db.Pool)
	auditRepo := NewAuditRepository(db.Pool)
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

	// Create audit log (success)
	log := &model.AuditLog{
		ID:        uuid.New().String(),
		UserID:    &user.ID,
		RequestID: "req-123",
		Action:    "command.process",
		Status:    "success",
		Metadata:  map[string]interface{}{"transcript": "open youtube"},
		CreatedAt: now,
	}
	if err := auditRepo.Create(ctx, log); err != nil {
		t.Fatalf("Create audit log failed: %v", err)
	}

	// Create audit log (error)
	errChain := "handler.ProcessCommand → service.ParseIntent: empty transcript"
	ip := "192.168.1.1"
	ua := "Mozilla/5.0"
	errorLog := &model.AuditLog{
		ID:         uuid.New().String(),
		UserID:     &user.ID,
		RequestID:  "req-456",
		Action:     "command.process",
		Status:     "error",
		ErrorChain: &errChain,
		IPAddress:  &ip,
		UserAgent:  &ua,
		Metadata:   map[string]interface{}{"transcript": ""},
		CreatedAt:  now,
	}
	if err := auditRepo.Create(ctx, errorLog); err != nil {
		t.Fatalf("Create error audit log failed: %v", err)
	}
}

func TestAuditRepository_NullUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	auditRepo := NewAuditRepository(db.Pool)
	ctx := context.Background()

	// Audit log without user (e.g., failed auth attempt)
	now := time.Now().Truncate(time.Microsecond)
	log := &model.AuditLog{
		ID:        uuid.New().String(),
		RequestID: "req-anon",
		Action:    "auth.login",
		Status:    "error",
		Metadata:  map[string]interface{}{"email": "unknown@test.com"},
		CreatedAt: now,
	}
	if err := auditRepo.Create(ctx, log); err != nil {
		t.Fatalf("Create audit log without user failed: %v", err)
	}
}
