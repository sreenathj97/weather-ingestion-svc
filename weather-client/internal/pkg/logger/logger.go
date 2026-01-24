package logger

import (
	"log/slog"
	"os"
)

// Logger is the global logger used across the application.
// Default to slog.Default() so callers don't panic if Init() isn't invoked.
var Logger = slog.Default()

// Init initializes the logging configuration
func Init() {
	Logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}
