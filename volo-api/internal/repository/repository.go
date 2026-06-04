package repository

import "github.com/jackc/pgx/v5/pgxpool"

// Repositories holds all database repositories.
type Repositories struct {
	User     *UserRepository
	Command  *CommandRepository
	Session  *SessionRepository
	Audit    *AuditRepository
	Settings *SettingsRepository
	Token    *TokenRepository
	Macro    *MacroRepository
}

// New creates all repositories with the given connection pool.
func New(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:     NewUserRepository(pool),
		Command:  NewCommandRepository(pool),
		Session:  NewSessionRepository(pool),
		Audit:    NewAuditRepository(pool),
		Settings: NewSettingsRepository(pool),
		Token:    NewTokenRepository(pool),
		Macro:    NewMacroRepository(pool),
	}
}
