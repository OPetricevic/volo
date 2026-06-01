package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/volo/volo-api/internal/intent"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/repository"
)

type CommandService struct {
	repos  *repository.Repositories
	rdb    *redis.Client
	parser *intent.HybridParser
}

func NewCommandService(repos *repository.Repositories, rdb *redis.Client) *CommandService {
	// Model directory — configurable via env, defaults to ./models
	modelDir := os.Getenv("VOLO_MODEL_PATH")
	if modelDir == "" {
		modelDir = "./models"
	}

	return &CommandService{
		repos:  repos,
		rdb:    rdb,
		parser: intent.NewHybridParser(modelDir),
	}
}

// Process parses a voice command and stores it in history.
func (s *CommandService) Process(ctx context.Context, userID string, req model.ProcessCommandRequest) (*model.CommandResponse, error) {
	if req.Transcript == "" {
		return nil, fmt.Errorf("service.Command.Process: transcript cannot be empty")
	}

	// Parse intent
	result := s.parser.Parse(req.Transcript)

	// Store in history
	cmd := &model.Command{
		ID:           uuid.New().String(),
		UserID:       userID,
		Transcript:   req.Transcript,
		ParsedAction: result.Action,
		ParsedTarget: result.Target,
		ParsedQuery:  result.Query,
		Confidence:   result.Confidence,
		ExecutedAt:   time.Now(),
	}
	if err := s.repos.Command.Create(ctx, cmd); err != nil {
		return nil, fmt.Errorf("service.Command.Process → Create: %w", err)
	}

	// Build response
	response := &model.CommandResponse{
		Action:     result.Action,
		Confidence: result.Confidence,
	}
	if result.Target != nil {
		response.Target = *result.Target
	}
	if result.Query != nil {
		response.Query = *result.Query
	}
	if result.ExecuteURL != "" {
		response.Execute = &model.ExecuteAction{
			URL:      result.ExecuteURL,
			AutoPlay: result.Action == "open-and-play",
		}
	}

	return response, nil
}

// GetHistory returns paginated command history for a user.
func (s *CommandService) GetHistory(ctx context.Context, userID string, page, pageSize int) (*model.HistoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	commands, total, err := s.repos.Command.GetByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("service.Command.GetHistory: %w", err)
	}

	return &model.HistoryResponse{
		Commands: commands,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// ClearHistory deletes all command history for a user.
func (s *CommandService) ClearHistory(ctx context.Context, userID string) error {
	if err := s.repos.Command.DeleteByUser(ctx, userID); err != nil {
		return fmt.Errorf("service.Command.ClearHistory: %w", err)
	}
	return nil
}

// GetSettings retrieves user settings.
func (s *CommandService) GetSettings(ctx context.Context, userID string) (*repository.UserSettings, error) {
	settings, err := s.repos.Settings.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Command.GetSettings: %w", err)
	}
	return settings, nil
}

// UpdateSettings updates user settings (partial update — only non-nil fields).
func (s *CommandService) UpdateSettings(ctx context.Context, userID string, req model.UpdateSettingsRequest) (*repository.UserSettings, error) {
	// Get current settings
	current, err := s.repos.Settings.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Command.UpdateSettings → Get: %w", err)
	}

	// Apply partial updates
	if req.MicMode != nil {
		current.MicMode = *req.MicMode
	}
	if req.WakeWord != nil {
		current.WakeWord = *req.WakeWord
	}
	if req.Language != nil {
		current.Language = *req.Language
	}

	// Save
	if err := s.repos.Settings.Update(ctx, userID, current); err != nil {
		return nil, fmt.Errorf("service.Command.UpdateSettings → Update: %w", err)
	}

	return current, nil
}

// GetSuggestions returns ranked suggestions based on user's command history.
// Uses frequency + recency weighting to surface the most likely next commands.
func (s *CommandService) GetSuggestions(ctx context.Context, userID string) (*model.SuggestionsResponse, error) {
	// Get recent commands (last 50)
	commands, _, err := s.repos.Command.GetByUser(ctx, userID, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("service.Command.GetSuggestions: %w", err)
	}

	if len(commands) == 0 {
		return &model.SuggestionsResponse{Suggestions: []model.Suggestion{}}, nil
	}

	// Build frequency map with recency weighting
	type scored struct {
		text      string
		category  string
		frequency int
		recency   float64 // higher = more recent
	}

	seen := make(map[string]*scored)
	now := time.Now()

	for i, cmd := range commands {
		key := cmd.ParsedAction
		if cmd.ParsedTarget != nil {
			key += ":" + *cmd.ParsedTarget
		}
		if cmd.ParsedQuery != nil {
			key += ":" + *cmd.ParsedQuery
		}

		if existing, ok := seen[key]; ok {
			existing.frequency++
			// Recency: more recent commands get higher weight
			age := now.Sub(cmd.ExecutedAt).Hours()
			recencyScore := 1.0 / (1.0 + age/24.0) // decay over days
			if recencyScore > existing.recency {
				existing.recency = recencyScore
			}
		} else {
			text := cmd.Transcript
			category := cmd.ParsedAction
			age := now.Sub(cmd.ExecutedAt).Hours()
			recencyScore := 1.0 / (1.0 + age/24.0)

			seen[key] = &scored{
				text:      text,
				category:  category,
				frequency: 1,
				recency:   recencyScore,
			}
			_ = i
		}
	}

	// Score and rank
	var suggestions []model.Suggestion
	for _, s := range seen {
		score := float64(s.frequency) * s.recency
		suggestions = append(suggestions, model.Suggestion{
			Text:  s.text,
			Type:  s.category,
			Score: score,
		})
	}

	// Sort by score descending
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Score > suggestions[i].Score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	// Cap at 10
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return &model.SuggestionsResponse{Suggestions: suggestions}, nil
}

// GetHistoryContext returns a formatted text summary of recent commands for LLM context.
func (s *CommandService) GetHistoryContext(ctx context.Context, userID string) (string, error) {
	commands, _, err := s.repos.Command.GetByUser(ctx, userID, 20, 0)
	if err != nil {
		return "", fmt.Errorf("service.Command.GetHistoryContext: %w", err)
	}

	if len(commands) == 0 {
		return "(No commands in history yet)", nil
	}

	var result string
	for _, cmd := range commands {
		target := ""
		if cmd.ParsedTarget != nil {
			target = " → " + *cmd.ParsedTarget
		}
		query := ""
		if cmd.ParsedQuery != nil {
			query = ": " + *cmd.ParsedQuery
		}
		result += fmt.Sprintf("- [%s] %s%s%s (%s)\n",
			cmd.ExecutedAt.Format("Jan 2 15:04"),
			cmd.ParsedAction,
			target,
			query,
			cmd.Transcript,
		)
	}
	return result, nil
}
