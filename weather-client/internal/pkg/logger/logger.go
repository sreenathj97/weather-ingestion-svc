package logger

import (
	"log/slog"
	"os"
)

// Logger is the global logger used across the application
var Logger *slog.Logger

// Init initializes the logging configuration
func Init() {
	Logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}
