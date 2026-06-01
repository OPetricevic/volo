package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/volo/volo-api/internal/config"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo represents the response from Google's userinfo endpoint.
type GoogleUserInfo struct {
	Sub     string `json:"sub"` // Unique Google user ID
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GoogleAuth handles the Google OAuth flow.
// 1. Exchange authorization code for tokens
// 2. Fetch user info from Google
// 3. Find or create user + credential
// 4. Link device
// 5. Issue JWT
func (s *AuthService) GoogleAuth(ctx context.Context, req model.GoogleAuthRequest) (*model.AuthResponse, error) {
	if req.Code == "" {
		return nil, fmt.Errorf("service.Auth.GoogleAuth: authorization code is required")
	}
	if req.DeviceID == "" {
		return nil, fmt.Errorf("service.Auth.GoogleAuth: device_id is required")
	}

	// 1. Build OAuth2 config
	oauthConfig := &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "postmessage", // For popup-based flow
		Scopes:       []string{"openid", "email", "profile"},
	}

	// 2. Exchange code for token
	token, err := oauthConfig.Exchange(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.GoogleAuth → Exchange: %w", err)
	}

	// 3. Fetch user info
	userInfo, err := fetchGoogleUserInfo(ctx, oauthConfig, token)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.GoogleAuth → fetchUserInfo: %w", err)
	}

	// 4. Find existing credential for this Google user
	cred, _ := s.repos.User.GetCredentialByProvider(ctx, "google", userInfo.Sub)

	now := time.Now()
	var user *model.User

	if cred != nil {
		// Existing user — just log them in
		user, err = s.repos.User.GetByID(ctx, cred.UserID)
		if err != nil {
			return nil, fmt.Errorf("service.Auth.GoogleAuth → GetByID: %w", err)
		}
	} else {
		// New user — check if email already exists (anonymous upgrade)
		existing, _ := s.repos.User.GetByEmail(ctx, userInfo.Email)
		if existing != nil {
			// Upgrade existing user with Google credential
			user = existing
		} else {
			// Create new user
			user = &model.User{
				ID:          uuid.New().String(),
				Email:       &userInfo.Email,
				DisplayName: &userInfo.Name,
				AvatarURL:   &userInfo.Picture,
				Role:        "user",
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			if err := s.repos.User.Create(ctx, user); err != nil {
				return nil, fmt.Errorf("service.Auth.GoogleAuth → CreateUser: %w", err)
			}
		}

		// Create Google credential
		providerID := userInfo.Sub
		newCred := &model.Credential{
			ID:             uuid.New().String(),
			UserID:         user.ID,
			Provider:       "google",
			ProviderUserID: &providerID,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := s.repos.User.CreateCredential(ctx, newCred); err != nil {
			return nil, fmt.Errorf("service.Auth.GoogleAuth → CreateCredential: %w", err)
		}
	}

	// 5. Link device (or update last seen)
	device := s.linkDevice(ctx, user.ID, req.DeviceID, now)

	// 6. Issue JWT
	jwtToken, err := s.issueToken(ctx, user.ID, &device.ID)
	if err != nil {
		return nil, fmt.Errorf("service.Auth.GoogleAuth → issueToken: %w", err)
	}

	return &model.AuthResponse{Token: jwtToken, User: *user, Device: *device}, nil
}

// linkDevice finds or creates a device for the user.
func (s *AuthService) linkDevice(ctx context.Context, userID, deviceID string, now time.Time) *model.Device {
	existing, _ := s.repos.User.GetDeviceByDeviceID(ctx, deviceID)
	if existing != nil {
		_ = s.repos.User.UpdateDeviceLastSeen(ctx, existing.ID)
		return existing
	}

	deviceName := "Unknown"
	device := &model.Device{
		ID:         uuid.New().String(),
		UserID:     userID,
		DeviceID:   deviceID,
		DeviceName: &deviceName,
		LastSeenAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_ = s.repos.User.CreateDevice(ctx, device)
	return device
}

// fetchGoogleUserInfo calls Google's userinfo endpoint.
func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google returned %d: %s", resp.StatusCode, string(body))
	}

	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	if info.Sub == "" || info.Email == "" {
		return nil, fmt.Errorf("incomplete user info from Google")
	}

	return &info, nil
}

// Ensure config is available (compile-time check)
var _ *config.Config
var _ *repository.Repositories
