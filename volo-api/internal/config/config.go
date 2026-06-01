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
}

func Load() *Config {
	return &Config{
		Port:               getEnv("VOLO_PORT", "8080"),
		DatabaseURL:        getEnv("VOLO_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/volo?sslmode=disable"),
		RedisAddr:          getEnv("VOLO_REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("VOLO_REDIS_PASSWORD", ""),
		RedisDB:            0,
		JWTSecret:          getEnv("VOLO_JWT_SECRET", "dev-secret-change-in-production"),
		CORSOrigins:        splitEnv("VOLO_CORS_ORIGINS", "*"),
		LogLevel:           getEnv("VOLO_LOG_LEVEL", "info"),
		GoogleClientID:     getEnv("VOLO_GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("VOLO_GOOGLE_CLIENT_SECRET", ""),
		SentryDSN:          getEnv("VOLO_SENTRY_DSN", ""),
		Environment:        getEnv("VOLO_ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitEnv(key, fallback string) []string {
	v := getEnv(key, fallback)
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
