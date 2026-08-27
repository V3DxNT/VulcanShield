package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type requestIDKey struct{}

// RequestIDKey is exported for handlers that need to read the request ID from context.
var RequestIDKey = requestIDKey{}

// RequestID is middleware that ensures every request has an X-Request-ID header.
// If the incoming request already carries X-Request-ID, that value is preserved.
// Otherwise a new random 16-byte hex ID is generated.
// The ID is echoed on the response and stored in the request context.
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = newRequestID()
			}
			w.Header().Set("X-Request-ID", id)
			ctx := context.WithValue(r.Context(), RequestIDKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID retrieves the request ID from a context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
