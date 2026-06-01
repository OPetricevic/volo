package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/volo/volo-api/internal/middleware"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/repository"
)

const (
	ollamaURL   = "http://ollama:11434/api/generate"
	gemmaModel  = "gemma2:2b"
	chatTimeout = 30 * time.Second
)

type ChatService struct {
	repos *repository.Repositories
}

func NewChatService(repos *repository.Repositories) *ChatService {
	return &ChatService{repos: repos}
}

type ChatRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id,omitempty"`
}

type ChatResponse struct {
	Reply   string       `json:"reply"`
	Sources []ChatSource `json:"sources,omitempty"`
}

type ChatSource struct {
	ID         string `json:"id"`
	Transcript string `json:"transcript"`
	ExecutedAt string `json:"executed_at"`
}

// Process handles a chat message by sending it to Ollama with history context.
func (s *ChatService) Process(ctx context.Context, userID string, req ChatRequest) (*ChatResponse, error) {
	if req.Message == "" {
		return nil, fmt.Errorf("service.Chat.Process: message cannot be empty")
	}

	// 1. Fetch user's recent history for context
	commands, _, err := s.repos.Command.GetByUser(ctx, userID, 20, 0)
	if err != nil {
		return nil, fmt.Errorf("service.Chat.Process → GetHistory: %w", err)
	}

	// 2. Build context string from history
	historyContext := buildHistoryContext(commands)

	// 3. Build prompt for Ollama
	prompt := buildChatPrompt(req.Message, historyContext)

	// 4. Call Ollama
	reply, err := callOllama(ctx, prompt)
	if err != nil {
		middleware.RecordOllamaRequest(false)
		return nil, fmt.Errorf("service.Chat.Process → callOllama: %w", err)
	}
	middleware.RecordOllamaRequest(true)

	// 5. Build sources from relevant commands
	var sources []ChatSource
	for _, cmd := range commands {
		sources = append(sources, ChatSource{
			ID:         cmd.ID,
			Transcript: cmd.Transcript,
			ExecutedAt: cmd.ExecutedAt.Format(time.RFC3339),
		})
	}

	return &ChatResponse{
		Reply:   reply,
		Sources: sources,
	}, nil
}

// CheckOllamaStatus checks if Ollama is reachable.
func (s *ChatService) CheckOllamaStatus() bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://ollama:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func buildHistoryContext(commands []model.Command) string {
	if len(commands) == 0 {
		return "(No commands in history yet)"
	}

	var sb strings.Builder
	for _, cmd := range commands {
		target := ""
		if cmd.ParsedTarget != nil {
			target = " → " + *cmd.ParsedTarget
		}
		query := ""
		if cmd.ParsedQuery != nil {
			query = ": " + *cmd.ParsedQuery
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s%s%s (%s)\n",
			cmd.ExecutedAt.Format("Jan 2 15:04"),
			cmd.ParsedAction,
			target,
			query,
			cmd.Transcript,
		))
	}
	return sb.String()
}

func buildChatPrompt(message, historyContext string) string {
	return fmt.Sprintf(`You are Volo, a voice assistant's history helper. The user will ask questions about their browsing and search history.

Given the user's question and their recent command history below, answer naturally and concisely. Keep responses short and helpful.

If the user asks about something not in their history, say so politely.
If the user gives a command (like "open youtube"), tell them to use voice mode instead.

Recent command history:
%s

User: %s

Volo:`, historyContext, message)
}

func callOllama(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, chatTimeout)
	defer cancel()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":  gemmaModel,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.3,
			"num_predict": 256,
		},
	})

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode ollama response: %w", err)
	}

	return result.Response, nil
}

// GenerateConversationID creates a new conversation ID.
func GenerateConversationID() string {
	return uuid.New().String()
}
