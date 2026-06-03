package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/volo/volo-api/internal/config"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/repository"
)

const tokenExpiry = 30 * 24 * time.Hour // 30 days

type AuthService struct {
	repos *repository.Repositories
	rdb   *redis.Client
	cfg   *config.Config
}

func NewAuthService(repos *repository.Repositories, rdb *redis.Client, cfg *config.Config) *AuthService {
	return &AuthService{repos: repos, rdb: rdb, cfg: cfg}
}

// RegisterDevice creates an anonymous user and device, returns a JWT.
func (s *AuthService) RegisterDevice(ctx context.Context, req model.RegisterDeviceRequest) (*model.AuthResponse, error) {
	now := time.Now()

	// Check if device already exists
	existing, _ := s.repos.User.GetDeviceByDeviceID(ctx, req.DeviceID)
	if existing != nil {
		// Device already registered — issue new token
		user, err := s.repos.User.GetByID(ctx, existing.UserID)
		if err != nil {
			return nil, fmt.Errorf("service.Auth.RegisterDevice → GetByID: %w", err)
		}
		token, err := s.issueToken(ctx, user.ID, &existing.ID)
		if err != nil {
			return nil, fmt.Errorf("service.Auth.RegisterDevice → issueToken: %w", err)
		}
		_ = s.repos.User.UpdateDeviceLastSeen(ctx, existing.ID)
		return &model.AuthResponse{Token: token, User: *user, Device: *existing}, nil
	}

	// Create new anonymous user
	user := &model.User{
		ID:        uuid.New().String(),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repos.User.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterDevice → CreateUser: %w", err)
	}

	// Create device
	deviceName := req.DeviceName
	platform := req.Platform
	device := &model.Device{
		ID:         uuid.New().String(),
		UserID:     user.ID,
		DeviceID:   req.DeviceID,
		DeviceName: &deviceName,
		Platform:   &platform,
		LastSeenAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.repos.User.CreateDevice(ctx, device); err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterDevice → CreateDevice: %w", err)
	}

	// Issue token
	token, err := s.issueToken(ctx, user.ID, &device.ID)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterDevice → issueToken: %w", err)
	}

	return &model.AuthResponse{Token: token, User: *user, Device: *device}, nil
}

// RegisterEmail upgrades an anonymous user or creates a new account with email+password.
func (s *AuthService) RegisterEmail(ctx context.Context, req model.RegisterEmailRequest, userID string) (*model.AuthResponse, error) {
	now := time.Now()

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterEmail → hash: %w", err)
	}

	// If user is already authenticated (upgrading anonymous), update their email
	if userID != "" {
		if err := s.repos.User.UpdateEmail(ctx, userID, req.Email); err != nil {
			return nil, fmt.Errorf("service.Auth.RegisterEmail → UpdateEmail: %w", err)
		}
		hashStr := string(hash)
		cred := &model.Credential{
			ID:           uuid.New().String(),
			UserID:       userID,
			Provider:     "password",
			PasswordHash: &hashStr,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := s.repos.User.CreateCredential(ctx, cred); err != nil {
			return nil, fmt.Errorf("service.Auth.RegisterEmail → CreateCredential: %w", err)
		}
		user, err := s.repos.User.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("service.Auth.RegisterEmail → GetByID: %w", err)
		}
		token, err := s.issueToken(ctx, userID, nil)
		if err != nil {
			return nil, fmt.Errorf("service.Auth.RegisterEmail → issueToken: %w", err)
		}
		device, _ := s.repos.User.GetDeviceByDeviceID(ctx, req.DeviceID)
		resp := &model.AuthResponse{Token: token, User: *user}
		if device != nil {
			resp.Device = *device
		}
		return resp, nil
	}

	// New user
	user := &model.User{
		ID:        uuid.New().String(),
		Email:     &req.Email,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repos.User.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterEmail → CreateUser: %w", err)
	}

	hashStr := string(hash)
	cred := &model.Credential{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		Provider:     "password",
		PasswordHash: &hashStr,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repos.User.CreateCredential(ctx, cred); err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterEmail → CreateCredential: %w", err)
	}

	// Link device
	deviceName := "Unknown"
	device := &model.Device{
		ID:         uuid.New().String(),
		UserID:     user.ID,
		DeviceID:   req.DeviceID,
		DeviceName: &deviceName,
		LastSeenAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.repos.User.CreateDevice(ctx, device); err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterEmail → CreateDevice: %w", err)
	}

	token, err := s.issueToken(ctx, user.ID, &device.ID)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.RegisterEmail → issueToken: %w", err)
	}
	return &model.AuthResponse{Token: token, User: *user, Device: *device}, nil
}

// LoginEmail authenticates with email+password.
func (s *AuthService) LoginEmail(ctx context.Context, req model.LoginEmailRequest) (*model.AuthResponse, error) {
	user, err := s.repos.User.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.LoginEmail → GetByEmail: invalid credentials")
	}

	cred, err := s.repos.User.GetCredentialByUserAndProvider(ctx, user.ID, "password")
	if err != nil || cred.PasswordHash == nil {
		return nil, fmt.Errorf("service.Auth.LoginEmail → GetCredential: invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*cred.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("service.Auth.LoginEmail → CompareHash: invalid credentials")
	}

	// Link device if new
	device, _ := s.repos.User.GetDeviceByDeviceID(ctx, req.DeviceID)
	if device == nil {
		now := time.Now()
		deviceName := "Unknown"
		device = &model.Device{
			ID:         uuid.New().String(),
			UserID:     user.ID,
			DeviceID:   req.DeviceID,
			DeviceName: &deviceName,
			LastSeenAt: now,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		_ = s.repos.User.CreateDevice(ctx, device)
	} else {
		_ = s.repos.User.UpdateDeviceLastSeen(ctx, device.ID)
	}

	token, err := s.issueToken(ctx, user.ID, &device.ID)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.LoginEmail → issueToken: %w", err)
	}

	return &model.AuthResponse{Token: token, User: *user, Device: *device}, nil
}

// ValidateToken checks if a JWT is valid and not revoked.
func (s *AuthService) ValidateToken(ctx context.Context, tokenStr string) (string, error) {
	// Parse JWT
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("service.Auth.ValidateToken: invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("service.Auth.ValidateToken: invalid claims")
	}

	userID, _ := claims["sub"].(string)
	if userID == "" {
		return "", fmt.Errorf("service.Auth.ValidateToken: missing sub claim")
	}

	// Check revocation via token hash
	hash := hashToken(tokenStr)

	// Check Redis cache first
	cached, err := s.rdb.Get(ctx, "session:"+hash).Result()
	if err == nil && cached == userID {
		return userID, nil
	}

	// Check database
	session, err := s.repos.Session.GetByTokenHash(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("service.Auth.ValidateToken: session not found or revoked")
	}

	// Cache in Redis (TTL = remaining session time)
	ttl := time.Until(session.ExpiresAt)
	if ttl > 0 {
		s.rdb.Set(ctx, "session:"+hash, userID, ttl)
	}

	return userID, nil
}

// Logout revokes the current session.
func (s *AuthService) Logout(ctx context.Context, tokenStr string) error {
	hash := hashToken(tokenStr)
	if err := s.repos.Session.Revoke(ctx, hash); err != nil {
		return fmt.Errorf("service.Auth.Logout: %w", err)
	}
	s.rdb.Del(ctx, "session:"+hash)
	return nil
}

// LogoutAll revokes all sessions for a user.
func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	if err := s.repos.Session.RevokeAllForUser(ctx, userID); err != nil {
		return fmt.Errorf("service.Auth.LogoutAll: %w", err)
	}
	// Note: Redis cache entries will expire naturally
	return nil
}

// GetUserByID retrieves a user by their ID.
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.GetUserByID: %w", err)
	}
	return user, nil
}

// GetUserDevice retrieves a device only if it belongs to the specified user.
func (s *AuthService) GetUserDevice(ctx context.Context, userID, deviceID string) (*model.Device, error) {
	device, err := s.repos.User.GetDeviceByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.GetUserDevice: %w", err)
	}
	if device == nil || device.UserID != userID {
		return nil, nil
	}
	return device, nil
}

// UnlinkDevice removes a device from a user's account.
func (s *AuthService) UnlinkDevice(ctx context.Context, userID, deviceID string) error {
	if err := s.repos.User.DeleteDevice(ctx, userID, deviceID); err != nil {
		return fmt.Errorf("service.Auth.UnlinkDevice: %w", err)
	}
	return nil
}

func (s *AuthService) issueToken(ctx context.Context, userID string, deviceID *string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(tokenExpiry)

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("service.Auth.issueToken → sign: %w", err)
	}

	// Store session
	session := &model.Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		DeviceID:  deviceID,
		TokenHash: hashToken(tokenStr),
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}
	if err := s.repos.Session.Create(ctx, session); err != nil {
		return "", fmt.Errorf("service.Auth.issueToken → CreateSession: %w", err)
	}

	return tokenStr, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// ForgotPassword generates a reset token and returns it (caller sends the email).
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (rawToken string, err error) {
	// Find user by email
	user, err := s.repos.User.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal whether email exists — return nil error
		slog.Info("forgot password for unknown email", "email", email)
		return "", nil
	}

	// Rate limit: max 3 resets per hour per user
	count, err := s.repos.Token.CountRecentByUserAndType(ctx, user.ID, "password_reset", time.Now().Add(-1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("service.Auth.ForgotPassword → CountRecent: %w", err)
	}
	if count >= 3 {
		slog.Warn("password reset rate limited", "user_id", user.ID)
		return "", nil // Silent — don't reveal rate limiting to the client
	}

	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("service.Auth.ForgotPassword → rand: %w", err)
	}
	rawToken = base64.URLEncoding.EncodeToString(tokenBytes)

	// Store hashed token
	token := &model.Token{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Type:      "password_reset",
		TokenHash: hashToken(rawToken),
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := s.repos.Token.Create(ctx, token); err != nil {
		return "", fmt.Errorf("service.Auth.ForgotPassword → Create: %w", err)
	}

	return rawToken, nil
}

// ResetPassword validates a reset token and updates the password.
func (s *AuthService) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	// Find token
	tokenHash := hashToken(rawToken)
	token, err := s.repos.Token.GetByHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("service.Auth.ResetPassword: invalid or expired token")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service.Auth.ResetPassword → hash: %w", err)
	}

	// Update credential
	hashStr := string(hash)
	if err := s.repos.User.UpdatePasswordHash(ctx, token.UserID, hashStr); err != nil {
		return fmt.Errorf("service.Auth.ResetPassword → UpdatePassword: %w", err)
	}

	// Mark token as used
	if err := s.repos.Token.MarkUsed(ctx, token.ID); err != nil {
		return fmt.Errorf("service.Auth.ResetPassword → MarkUsed: %w", err)
	}

	// Invalidate all other reset tokens for this user
	_ = s.repos.Token.InvalidateAllForUser(ctx, token.UserID, "password_reset")

	// Revoke all sessions (force re-login everywhere)
	_ = s.repos.Session.RevokeAllForUser(ctx, token.UserID)

	slog.Info("password reset successful", "user_id", token.UserID)
	return nil
}

// SendVerificationEmail generates an email verification token and returns it.
func (s *AuthService) SendVerificationEmail(ctx context.Context, userID, email string) (string, error) {
	// Rate limit: max 3 per hour
	count, err := s.repos.Token.CountRecentByUserAndType(ctx, userID, "email_verify", time.Now().Add(-1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("service.Auth.SendVerificationEmail → CountRecent: %w", err)
	}
	if count >= 3 {
		return "", nil // Silent rate limit
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("service.Auth.SendVerificationEmail → rand: %w", err)
	}
	rawToken := base64.URLEncoding.EncodeToString(tokenBytes)

	token := &model.Token{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      "email_verify",
		TokenHash: hashToken(rawToken),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours for email verification
		CreatedAt: time.Now(),
	}
	if err := s.repos.Token.Create(ctx, token); err != nil {
		return "", fmt.Errorf("service.Auth.SendVerificationEmail → Create: %w", err)
	}

	return rawToken, nil
}

// VerifyEmail validates a verification token and marks the user's email as verified.
func (s *AuthService) VerifyEmail(ctx context.Context, rawToken string) error {
	tokenHash := hashToken(rawToken)
	token, err := s.repos.Token.GetByHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("service.Auth.VerifyEmail: invalid or expired token")
	}
	if token.Type != "email_verify" {
		return fmt.Errorf("service.Auth.VerifyEmail: invalid token type")
	}

	// Mark email as verified
	if err := s.repos.User.VerifyEmail(ctx, token.UserID); err != nil {
		return fmt.Errorf("service.Auth.VerifyEmail → VerifyEmail: %w", err)
	}

	// Mark token as used
	if err := s.repos.Token.MarkUsed(ctx, token.ID); err != nil {
		return fmt.Errorf("service.Auth.VerifyEmail → MarkUsed: %w", err)
	}

	// Invalidate other verification tokens for this user
	_ = s.repos.Token.InvalidateAllForUser(ctx, token.UserID, "email_verify")

	slog.Info("email verified", "user_id", token.UserID)
	return nil
}

// ChangePassword validates the current password and updates to the new one.
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Get current credential
	cred, err := s.repos.User.GetCredentialByUserAndProvider(ctx, userID, "password")
	if err != nil || cred.PasswordHash == nil {
		return fmt.Errorf("service.Auth.ChangePassword: no password credential found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(*cred.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("service.Auth.ChangePassword: current password is incorrect")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service.Auth.ChangePassword → hash: %w", err)
	}

	// Update
	if err := s.repos.User.UpdatePasswordHash(ctx, userID, string(hash)); err != nil {
		return fmt.Errorf("service.Auth.ChangePassword → Update: %w", err)
	}

	slog.Info("password changed", "user_id", userID)
	return nil
}

// DeleteAccount soft-deletes the user after verifying their password.
func (s *AuthService) DeleteAccount(ctx context.Context, userID, password string) error {
	// Verify password (require confirmation for destructive action)
	cred, err := s.repos.User.GetCredentialByUserAndProvider(ctx, userID, "password")
	if err != nil || cred.PasswordHash == nil {
		return fmt.Errorf("service.Auth.DeleteAccount: no password credential found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*cred.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("service.Auth.DeleteAccount: incorrect password")
	}

	// Soft delete user
	if err := s.repos.User.SoftDelete(ctx, userID); err != nil {
		return fmt.Errorf("service.Auth.DeleteAccount → SoftDelete: %w", err)
	}

	// Revoke all sessions
	_ = s.repos.Session.RevokeAllForUser(ctx, userID)

	slog.Info("account deleted", "user_id", userID)
	return nil
}
