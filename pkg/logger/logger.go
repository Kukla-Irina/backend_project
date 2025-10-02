package logger

import (
	"log/slog"
	"os"
)

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

func L() *slog.Logger {
	return defaultLogger
}
