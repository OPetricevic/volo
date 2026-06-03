package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/volo/volo-api/internal/config"
	"github.com/volo/volo-api/internal/handler"
	"github.com/volo/volo-api/internal/middleware"
	"github.com/volo/volo-api/internal/repository"
	"github.com/volo/volo-api/internal/service"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load config
	cfg := config.Load()

	// Logger
	logLevel := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Warn if using default JWT secret in production
	if cfg.Environment == "production" && cfg.JWTSecret == "dev-secret-change-in-production" {
		slog.Error("CRITICAL: Using default JWT secret in production. Set VOLO_JWT_SECRET env var.")
		os.Exit(1)
	}

	// Sentry
	middleware.InitSentry(cfg.SentryDSN, cfg.Environment)

	// Postgres
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to redis")

	// Repositories
	repos := repository.New(pool)

	// Services
	services := service.New(repos, rdb, cfg)

	// Handlers
	handlers := handler.New(services, cfg)

	// Router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Metrics)
	r.Use(middleware.SentryMiddleware())
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.ValidateBody)

	// Health
	r.Get("/health", handlers.Health)
	r.Get("/api/health", handlers.Health)

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Public (no auth)
		r.Post("/auth/device", handlers.RegisterDevice)
		r.Post("/auth/google", handlers.GoogleAuth)
		r.Post("/auth/register", handlers.RegisterEmail)
		r.Post("/auth/login", handlers.LoginEmail)
		r.Post("/auth/forgot-password", handlers.ForgotPassword)
		r.Post("/auth/reset-password", handlers.ResetPassword)
		r.Post("/auth/verify-email", handlers.VerifyEmail)

		// Protected (JWT required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(services.Auth, cfg.JWTSecret))
			r.Use(middleware.RateLimit(rdb, 60, time.Minute))

			r.Post("/command", handlers.ProcessCommand)
			r.Get("/suggestions", handlers.GetSuggestions)
			r.Get("/history", handlers.GetHistory)
			r.Delete("/history", handlers.ClearHistory)
			r.Get("/history/context", handlers.GetHistoryContext)

			r.Get("/settings", handlers.GetSettings)
			r.Put("/settings", handlers.UpdateSettings)

			r.Post("/chat", handlers.Chat)
			r.Get("/chat/status", handlers.ChatStatus)

			r.Post("/auth/logout", handlers.Logout)
			r.Post("/auth/logout-all", handlers.LogoutAll)
			r.Post("/auth/change-password", handlers.ChangePassword)
			r.Post("/auth/resend-verification", handlers.ResendVerification)
			r.Delete("/auth/account", handlers.DeleteAccount)
			r.Delete("/auth/device/{deviceID}", handlers.UnlinkDevice)
		})
	})

	// Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second, // Must exceed Ollama chat timeout (30s)
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}
