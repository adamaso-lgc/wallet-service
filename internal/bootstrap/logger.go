package bootstrap

import (
	"log/slog"
	"os"
)

// NewLogger creates a structured logger. JSON in production, human-readable text locally.
func NewLogger(env string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}

	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
