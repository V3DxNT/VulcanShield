package middleware

import "net/http"

// Chain composes multiple middleware into a single middleware.
// Middleware are applied in the order provided (first = outermost).
//
// Example:
//
//	Chain(RequestID(), Logger(), Recovery())(mux)
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		// Apply in reverse so the first middleware is outermost
		h := final
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
