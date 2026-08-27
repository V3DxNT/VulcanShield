package server

import (
	"context"
	"log/slog"
	"net/http"

	v1 "github.com/vulcanshield/backend/internal/api/v1"
	"github.com/vulcanshield/backend/internal/middleware"
)

// Dependencies holds all wired dependencies injected into the router.
// Each field is optional in the sense that nil values degrade gracefully,
// but PostgreSQL is required at startup per project architecture.
type Dependencies struct {
	Logger *slog.Logger
	Health *v1.HealthHandlers
}

// NewRouter builds the HTTP handler tree using stdlib http.ServeMux.
// Uses Go 1.22+ method-qualified routing patterns.
// Middleware applied (outer → inner): RequestID, StructuredLogger, Recovery.
// CORS is isolated in middleware/cors.go and NOT applied in Phase 3.
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	// Health routes — outside /api/v1 for load-balancer / container health check compatibility
	mux.HandleFunc("GET /health", deps.Health.Health)
	mux.HandleFunc("GET /ready", deps.Health.Ready)

	// API v1 routes — Phase 4+ will add routes here
	v1.Mount(mux, deps.Health)

	// Apply middleware chain (outermost first)
	return middleware.Chain(
		middleware.RequestID(),
		middleware.StructuredLogger(deps.Logger),
		middleware.Recovery(deps.Logger),
		// CORS intentionally omitted — enabled in Phase 12 (Next.js frontend)
	)(mux)
}

// pgxPinger wraps a pgxpool.Pool to satisfy the v1.Prober interface.
// Defined here to avoid circular imports between server and db packages.
type pgxPinger struct {
	ping func(ctx context.Context) error
}

func (p *pgxPinger) Ping(ctx context.Context) error {
	return p.ping(ctx)
}

// NewPgxProber wraps a ping function as a v1.Prober.
func NewPgxProber(ping func(ctx context.Context) error) v1.Prober {
	return &pgxPinger{ping: ping}
}
