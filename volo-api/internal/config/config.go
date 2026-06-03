package config

import "os"

type Config struct {
	Port               string
	DatabaseURL        string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	JWTSecret          string
	CORSOrigins        []string
	LogLevel           string
	GoogleClientID     string
	GoogleClientSecret string
	SentryDSN          string
	Environment        string
	ResendAPIKey       string
	ResetPasswordURL   string
	VerifyEmailURL     string
}

func Load() *Config {
	return &Config{
		Port:               getEnvAny([]string{"VOLO_PORT", "PORT"}, "8080"),
		DatabaseURL:        getEnvAny([]string{"VOLO_DATABASE_URL", "DATABASE_URL"}, "postgres://postgres:postgres@localhost:5432/volo?sslmode=disable"),
		RedisAddr:          getEnvAny([]string{"VOLO_REDIS_ADDR", "REDIS_ADDR"}, "localhost:6379"),
		RedisPassword:      getEnvAny([]string{"VOLO_REDIS_PASSWORD", "REDIS_PASSWORD"}, ""),
		RedisDB:            0,
		JWTSecret:          getEnvAny([]string{"VOLO_JWT_SECRET", "JWT_SECRET"}, "dev-secret-change-in-production"),
		CORSOrigins:        splitEnvAny([]string{"VOLO_CORS_ORIGINS", "CORS_ALLOWED_ORIGINS"}, "*"),
		LogLevel:           getEnvAny([]string{"VOLO_LOG_LEVEL", "LOG_LEVEL"}, "info"),
		GoogleClientID:     getEnvAny([]string{"VOLO_GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_ID"}, ""),
		GoogleClientSecret: getEnvAny([]string{"VOLO_GOOGLE_CLIENT_SECRET", "GOOGLE_CLIENT_SECRET"}, ""),
		SentryDSN:          getEnvAny([]string{"VOLO_SENTRY_DSN", "SENTRY_DSN"}, ""),
		Environment:        getEnvAny([]string{"VOLO_ENVIRONMENT", "APP_ENV"}, "development"),
		ResendAPIKey:       getEnvAny([]string{"VOLO_RESEND_API_KEY", "RESEND_API_KEY"}, ""),
		ResetPasswordURL:   getEnvAny([]string{"VOLO_RESET_PASSWORD_URL", "RESET_PASSWORD_URL"}, "http://localhost:5173/reset"),
		VerifyEmailURL:     getEnvAny([]string{"VOLO_VERIFY_EMAIL_URL", "VERIFY_EMAIL_URL"}, "http://localhost:5173/verify"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return fallback
}

func splitEnvAny(keys []string, fallback string) []string {
	v := getEnvAny(keys, fallback)
	if v == "*" {
		return []string{"*"}
	}
	var result []string
	current := ""
	for _, ch := range v {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func splitEnv(key, fallback string) []string {
	return splitEnvAny([]string{key}, fallback)
}
