package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a structured JSON logger that writes exclusively to stdout.
// Rotation and retention are deliberately delegated to the container runtime.
func New(level string) *slog.Logger {
	options := &slog.HandlerOptions{Level: parseLevel(level)}
	return slog.New(slog.NewJSONHandler(os.Stdout, options))
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
