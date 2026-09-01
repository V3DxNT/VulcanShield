package server

import (
	"context"
	"log/slog"
	"net/http"

	v1 "github.com/vulcanshield/backend/internal/api/v1"
	"github.com/vulcanshield/backend/internal/middleware"
)




type Dependencies struct {
	Logger   *slog.Logger
	Health   *v1.HealthHandlers
	Handlers *v1.Handlers 
}





func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	
	mux.HandleFunc("GET /health", deps.Health.Health)
	mux.HandleFunc("GET /ready", deps.Health.Ready)

	
	if deps.Handlers != nil {
		v1.Mount(mux, deps.Handlers)
	}

	
	return middleware.Chain(
		middleware.RequestID(),
		middleware.StructuredLogger(deps.Logger),
		middleware.Recovery(deps.Logger),
		
	)(mux)
}



type pgxPinger struct {
	ping func(ctx context.Context) error
}

func (p *pgxPinger) Ping(ctx context.Context) error {
	return p.ping(ctx)
}


func NewPgxProber(ping func(ctx context.Context) error) v1.Prober {
	return &pgxPinger{ping: ping}
}
