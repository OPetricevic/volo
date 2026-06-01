package middleware

import (
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// InitSentry initializes the Sentry SDK. Call once at startup.
func InitSentry(dsn string, environment string) {
	if dsn == "" {
		slog.Info("Sentry DSN not configured, error tracking disabled")
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		Release:          "volo-api@0.1.0",
		TracesSampleRate: 0.2, // 20% of requests get performance tracing
		EnableTracing:    true,
	})
	if err != nil {
		slog.Error("failed to initialize Sentry", "error", err)
		return
	}

	slog.Info("Sentry initialized", "environment", environment)
}

// SentryMiddleware wraps the Sentry HTTP handler for automatic error capture.
func SentryMiddleware() func(http.Handler) http.Handler {
	handler := sentryhttp.New(sentryhttp.Options{
		Repanic: true, // Let our Recoverer handle the panic after Sentry captures it
	})
	return func(next http.Handler) http.Handler {
		return handler.Handle(next)
	}
}

// CaptureError sends an error to Sentry with context.
func CaptureError(err error, tags map[string]string) {
	if err == nil {
		return
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		for k, v := range tags {
			scope.SetTag(k, v)
		}
		sentry.CaptureException(err)
	})
}

// CaptureMessage sends a message to Sentry.
func CaptureMessage(msg string, level sentry.Level) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
		sentry.CaptureMessage(msg)
	})
}
