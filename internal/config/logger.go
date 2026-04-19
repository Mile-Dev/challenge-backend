package config

import (
	"log/slog"
	"os"
)

// NewLogger crea un logger estructurado en formato JSON.
func NewLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")

	if env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
