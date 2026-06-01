package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/testutil"
)

func TestUserRepository_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	repo := NewUserRepository(db.Pool)
	ctx := context.Background()

	// Create user
	now := time.Now().Truncate(time.Microsecond)
	email := "test@example.com"
	user := &model.User{
		ID:        uuid.New().String(),
		Email:     &email,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Get by ID
	got, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("ID = %q, want %q", got.ID, user.ID)
	}
	if got.Email == nil || *got.Email != email {
		t.Errorf("Email = %v, want %q", got.Email, email)
	}
	if got.Role != "user" {
		t.Errorf("Role = %q, want user", got.Role)
	}

	// Get by email
	got2, err := repo.GetByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if got2.ID != user.ID {
		t.Errorf("GetByEmail returned wrong user: %q", got2.ID)
	}
}

func TestUserRepository_Credentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	repo := NewUserRepository(db.Pool)
	ctx := context.Background()

	// Create user first
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create credential
	hash := "$2a$10$fakehashvalue"
	cred := &model.Credential{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		Provider:     "password",
		PasswordHash: &hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := repo.CreateCredential(ctx, cred); err != nil {
		t.Fatalf("CreateCredential failed: %v", err)
	}

	// Get credential
	got, err := repo.GetCredentialByUserAndProvider(ctx, user.ID, "password")
	if err != nil {
		t.Fatalf("GetCredentialByUserAndProvider failed: %v", err)
	}
	if got.ID != cred.ID {
		t.Errorf("credential ID = %q, want %q", got.ID, cred.ID)
	}
	if got.PasswordHash == nil || *got.PasswordHash != hash {
		t.Error("password hash mismatch")
	}
}

func TestUserRepository_Devices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupPostgres(t)
	repo := NewUserRepository(db.Pool)
	ctx := context.Background()

	// Create user
	now := time.Now().Truncate(time.Microsecond)
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create user failed: %v", err)
	}

	// Create device
	deviceName := "Chrome on Windows"
	platform := "extension"
	device := &model.Device{
		ID:         uuid.New().String(),
		UserID:     user.ID,
		DeviceID:   "device-abc-123",
		DeviceName: &deviceName,
		Platform:   &platform,
		LastSeenAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := repo.CreateDevice(ctx, device); err != nil {
		t.Fatalf("CreateDevice failed: %v", err)
	}

	// Get by device ID
	got, err := repo.GetDeviceByDeviceID(ctx, "device-abc-123")
	if err != nil {
		t.Fatalf("GetDeviceByDeviceID failed: %v", err)
	}
	if got.ID != device.ID {
		t.Errorf("device ID = %q, want %q", got.ID, device.ID)
	}
	if got.UserID != user.ID {
		t.Errorf("device user_id = %q, want %q", got.UserID, user.ID)
	}

	// Update last seen
	if err := repo.UpdateDeviceLastSeen(ctx, device.ID); err != nil {
		t.Fatalf("UpdateDeviceLastSeen failed: %v", err)
	}

	// Soft delete
	if err := repo.SoftDeleteDevice(ctx, device.ID); err != nil {
		t.Fatalf("SoftDeleteDevice failed: %v", err)
	}

	// Should not find after soft delete
	_, err = repo.GetDeviceByDeviceID(ctx, "device-abc-123")
	if err == nil {
		t.Error("expected error after soft delete, got nil")
	}
}
