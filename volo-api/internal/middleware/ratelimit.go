package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/volo/volo-api/internal/model"
)

// RateLimit implements a sliding window rate limiter using Redis.
func RateLimit(rdb *redis.Client, maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r.Context())
			if userID == "" {
				next.ServeHTTP(w, r)
				return
			}

			key := fmt.Sprintf("ratelimit:%s", userID)
			ctx := r.Context()

			allowed, err := checkRateLimit(ctx, rdb, key, maxRequests, window)
			if err != nil {
				// If Redis is down, allow the request (fail open)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				RecordRateLimitHit()
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(model.Response{
					Error: &model.ErrorResponse{
						Code:    "RATE_LIMITED",
						Message: "Too many requests. Please slow down.",
					},
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func checkRateLimit(ctx context.Context, rdb *redis.Client, key string, max int, window time.Duration) (bool, error) {
	now := time.Now().UnixNano()
	windowStart := now - window.Nanoseconds()

	pipe := rdb.Pipeline()
	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	// Add current request (use nanoseconds as both score and member for uniqueness)
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
	// Count requests in window
	countCmd := pipe.ZCard(ctx, key)
	// Set expiry on the key
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := countCmd.Val()
	return count <= int64(max), nil
}
