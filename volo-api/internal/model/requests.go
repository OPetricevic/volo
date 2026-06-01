package model

// --- Auth Requests ---

type RegisterDeviceRequest struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Platform   string `json:"platform"` // "extension", "desktop"
}

type GoogleAuthRequest struct {
	Code     string `json:"code"`
	DeviceID string `json:"device_id"`
}

type RegisterEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

type LoginEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

// --- Command Requests ---

type ProcessCommandRequest struct {
	Transcript string         `json:"transcript"`
	Context    CommandContext `json:"context,omitempty"`
}

type CommandContext struct {
	CurrentURL string `json:"current_url,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}

// --- Settings ---

type UpdateSettingsRequest struct {
	MicMode  *string `json:"mic_mode,omitempty"`
	WakeWord *string `json:"wake_word,omitempty"`
	Language *string `json:"language,omitempty"`
}
