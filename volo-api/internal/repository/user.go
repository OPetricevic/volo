package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/volo/volo-api/internal/model"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, display_name, email, avatar_url, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.ID, user.DisplayName, user.Email, user.AvatarURL, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.User.Create: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, display_name, email, avatar_url, role, email_verified_at, created_at, updated_at
		 FROM users WHERE id = $1 AND deleted_at IS NULL`, id,
	).Scan(&user.ID, &user.DisplayName, &user.Email, &user.AvatarURL, &user.Role, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.User.GetByID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, display_name, email, avatar_url, role, email_verified_at, created_at, updated_at
		 FROM users WHERE email = $1 AND deleted_at IS NULL`, email,
	).Scan(&user.ID, &user.DisplayName, &user.Email, &user.AvatarURL, &user.Role, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.User.GetByEmail: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateEmail(ctx context.Context, id string, email string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET email = $1, updated_at = now() WHERE id = $2`, email, id,
	)
	if err != nil {
		return fmt.Errorf("repository.User.UpdateEmail: %w", err)
	}
	return nil
}

// --- Credentials ---

func (r *UserRepository) CreateCredential(ctx context.Context, cred *model.Credential) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO credentials (id, user_id, provider, provider_user_id, password_hash, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		cred.ID, cred.UserID, cred.Provider, cred.ProviderUserID, cred.PasswordHash, cred.CreatedAt, cred.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.User.CreateCredential: %w", err)
	}
	return nil
}

func (r *UserRepository) GetCredentialByProvider(ctx context.Context, provider, providerUserID string) (*model.Credential, error) {
	var cred model.Credential
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, provider, provider_user_id, password_hash, created_at, updated_at
		 FROM credentials WHERE provider = $1 AND provider_user_id = $2`, provider, providerUserID,
	).Scan(&cred.ID, &cred.UserID, &cred.Provider, &cred.ProviderUserID, &cred.PasswordHash, &cred.CreatedAt, &cred.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.User.GetCredentialByProvider: %w", err)
	}
	return &cred, nil
}

func (r *UserRepository) GetCredentialByUserAndProvider(ctx context.Context, userID, provider string) (*model.Credential, error) {
	var cred model.Credential
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, provider, provider_user_id, password_hash, created_at, updated_at
		 FROM credentials WHERE user_id = $1 AND provider = $2`, userID, provider,
	).Scan(&cred.ID, &cred.UserID, &cred.Provider, &cred.ProviderUserID, &cred.PasswordHash, &cred.CreatedAt, &cred.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.User.GetCredentialByUserAndProvider: %w", err)
	}
	return &cred, nil
}

// --- Devices ---

func (r *UserRepository) CreateDevice(ctx context.Context, device *model.Device) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO devices (id, user_id, device_id, device_name, platform, last_seen_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		device.ID, device.UserID, device.DeviceID, device.DeviceName, device.Platform, device.LastSeenAt, device.CreatedAt, device.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository.User.CreateDevice: %w", err)
	}
	return nil
}

func (r *UserRepository) GetDeviceByDeviceID(ctx context.Context, deviceID string) (*model.Device, error) {
	var device model.Device
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, device_id, device_name, platform, last_seen_at, created_at, updated_at
		 FROM devices WHERE device_id = $1 AND deleted_at IS NULL`, deviceID,
	).Scan(&device.ID, &device.UserID, &device.DeviceID, &device.DeviceName, &device.Platform, &device.LastSeenAt, &device.CreatedAt, &device.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.User.GetDeviceByDeviceID: %w", err)
	}
	return &device, nil
}

func (r *UserRepository) UpdateDeviceLastSeen(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE devices SET last_seen_at = now(), updated_at = now() WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("repository.User.UpdateDeviceLastSeen: %w", err)
	}
	return nil
}

func (r *UserRepository) SoftDeleteDevice(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE devices SET deleted_at = now(), updated_at = now() WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("repository.User.SoftDeleteDevice: %w", err)
	}
	return nil
}

// DeleteDevice soft-deletes a device for a specific user (ownership check).
func (r *UserRepository) DeleteDevice(ctx context.Context, userID, deviceID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE devices SET deleted_at = now(), updated_at = now()
		 WHERE device_id = $1 AND user_id = $2 AND deleted_at IS NULL`, deviceID, userID,
	)
	if err != nil {
		return fmt.Errorf("repository.User.DeleteDevice: %w", err)
	}
	return nil
}

// UpdatePasswordHash updates the password hash for a user's password credential.
func (r *UserRepository) UpdatePasswordHash(ctx context.Context, userID, passwordHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE credentials SET password_hash = $1, updated_at = now()
		 WHERE user_id = $2 AND provider = 'password'`, passwordHash, userID,
	)
	if err != nil {
		return fmt.Errorf("repository.User.UpdatePasswordHash: %w", err)
	}
	return nil
}

// VerifyEmail sets the email_verified_at timestamp.
func (r *UserRepository) VerifyEmail(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET email_verified_at = now(), updated_at = now() WHERE id = $1`, userID,
	)
	if err != nil {
		return fmt.Errorf("repository.User.VerifyEmail: %w", err)
	}
	return nil
}

// SoftDelete marks a user as deleted.
func (r *UserRepository) SoftDelete(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET deleted_at = now(), updated_at = now() WHERE id = $1`, userID,
	)
	if err != nil {
		return fmt.Errorf("repository.User.SoftDelete: %w", err)
	}
	return nil
}
