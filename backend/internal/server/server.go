package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Server wraps the stdlib http.Server with structured lifecycle management.
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

// New creates a Server ready to serve on the given port.
func New(port string, handler http.Handler, log *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         net.JoinHostPort("", port),
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		log: log,
	}
}

// Run starts the HTTP server. It blocks until the server stops.
// Returns nil on graceful shutdown, or an error for unexpected failures.
func (s *Server) Run() error {
	s.log.Info("http server starting", "addr", s.httpServer.Addr, "service", "backend")
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown gracefully drains in-flight requests within the provided context deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("http server shutting down", "service", "backend")
	return s.httpServer.Shutdown(ctx)
}
