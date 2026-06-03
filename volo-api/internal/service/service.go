package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/volo/volo-api/internal/config"
	"github.com/volo/volo-api/internal/repository"
)

// Services holds all service instances.
type Services struct {
	Auth    *AuthService
	Command *CommandService
	Chat    *ChatService
	Email   *EmailService
}

// New creates all services.
func New(repos *repository.Repositories, rdb *redis.Client, cfg *config.Config) *Services {
	auth := NewAuthService(repos, rdb, cfg)
	command := NewCommandService(repos, rdb)
	chat := NewChatService(repos)
	email := NewEmailService(cfg.ResendAPIKey)

	return &Services{
		Auth:    auth,
		Command: command,
		Chat:    chat,
		Email:   email,
	}
}
