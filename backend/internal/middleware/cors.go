package middleware

import (
	"net/http"
	"strings"
)

// CORS returns middleware that sets Cross-Origin Resource Sharing headers.
// This middleware is DISABLED BY DEFAULT in Phase 3 — there is no frontend yet.
// It will be enabled in Phase 12 when the Next.js frontend is introduced.
//
// Usage (when enabling):
//
//	origins := []string{"http://localhost:3000"}
//	middleware.Chain(CORS(origins), ...)(mux)
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[strings.TrimRight(o, "/")] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if originSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.Header().Set("Vary", "Origin")
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
