package logger

import (
	"log/slog"
	"os"
	"strings"

	"giggler-golang/src/shared/must"
)

func InitSlog() Interface {
	var slogLevel slog.Level

	switch strings.ToLower(must.GetEnv("GIGGLER_LOGGER_LEVEL")) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		panic(must.ErrInvalidEnvValue)
	}

	return slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slogLevel},
	))
}
