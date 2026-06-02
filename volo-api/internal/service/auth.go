package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
