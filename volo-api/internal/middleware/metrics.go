package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "volo_http_requests_total",
			Help: "Total HTTP requests by method, path, and status code",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "volo_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	commandsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "volo_commands_processed_total",
			Help: "Voice commands processed by action type",
		},
		[]string{"action"},
	)

	authRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "volo_auth_registrations_total",
			Help: "Total device registrations",
		},
	)

	rateLimitHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "volo_rate_limit_hits_total",
			Help: "Total rate limit rejections",
		},
	)

	ollamaRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "volo_ollama_requests_total",
			Help: "Requests to Ollama by status",
		},
		[]string{"status"},
	)
)

// Metrics middleware records request count and duration.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &metricsWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.statusCode)
		path := normalizePath(r.URL.Path)

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

// RecordCommand records a processed command metric.
func RecordCommand(action string) {
	commandsProcessed.WithLabelValues(action).Inc()
}

// RecordRegistration records a new device registration.
func RecordRegistration() {
	authRegistrations.Inc()
}

// RecordRateLimitHit records a rate limit rejection.
func RecordRateLimitHit() {
	rateLimitHits.Inc()
}

// RecordOllamaRequest records an Ollama request result.
func RecordOllamaRequest(success bool) {
	if success {
		ollamaRequests.WithLabelValues("success").Inc()
	} else {
		ollamaRequests.WithLabelValues("error").Inc()
	}
}

type metricsWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *metricsWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// normalizePath reduces cardinality by grouping dynamic path segments.
func normalizePath(path string) string {
	// Group all UUIDs and IDs to prevent metric explosion
	switch {
	case len(path) > 20:
		// Truncate very long paths
		return path[:20]
	default:
		return path
	}
}
