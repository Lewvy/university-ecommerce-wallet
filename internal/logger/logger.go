package logger

import (
	"log/slog"
	"os"
)

func NewLogger(env string) *slog.Logger {
	switch env {
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}
