package model

// --- Standard Response Envelope ---

type Response struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	InternalError string `json:"internal_error,omitempty"`
}

// --- Auth Responses ---

type AuthResponse struct {
	Token  string `json:"token"`
	User   User   `json:"user"`
	Device Device `json:"device"`
}

// --- Command Responses ---

type CommandResponse struct {
	Action      string         `json:"action"`
	Target      string         `json:"target,omitempty"`
	Query       string         `json:"query,omitempty"`
	Confidence  float64        `json:"confidence"`
	Suggestions []string       `json:"suggestions,omitempty"`
	Execute     *ExecuteAction `json:"execute,omitempty"`
}

type ExecuteAction struct {
	URL      string `json:"url"`
	AutoPlay bool   `json:"auto_play,omitempty"`
}

// --- Suggestions ---

type SuggestionsResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Text  string  `json:"text"`
	Type  string  `json:"type"` // "search", "site", "command"
	Score float64 `json:"score"`
}

// --- History ---

type HistoryResponse struct {
	Commands []Command `json:"commands"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
