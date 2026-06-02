package middleware

import "net/http"

// CORS handles Cross-Origin Resource Sharing headers.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if isAllowAll(allowedOrigins) {
				// Development mode: allow all origins but WITHOUT credentials
				// (browsers reject wildcard + credentials anyway)
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if isAllowed(origin, allowedOrigins) {
				// Production mode: reflect specific origin with credentials
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAllowAll(allowed []string) bool {
	return len(allowed) == 1 && allowed[0] == "*"
}

func isAllowed(origin string, allowed []string) bool {
	if origin == "" || len(allowed) == 0 {
		return false
	}
	for _, a := range allowed {
		if a == origin {
			return true
		}
	}
	return false
}
