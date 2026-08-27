package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey struct{}

// Init initialises a structured slog.Logger based on the given level string.
// JSON handler is used in all cases for machine-readable structured output.
// Sets the default logger so stdlib log calls also emit structured output.
func Init(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}

	var handler slog.Handler
	if level == "debug" {
		// Human-readable text for local development
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// JSON for production / Docker
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	log := slog.New(handler)
	slog.SetDefault(log)
	return log
}

// WithContext returns a new context carrying the given logger.
func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

// FromContext retrieves the logger stored in ctx.
// Falls back to the default slog logger if none is set.
func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(contextKey{}).(*slog.Logger); ok && log != nil {
		return log
	}
	return slog.Default()
}
