package log

import (
	"io"
	"log/slog"
	"os"

	"giggler-golang/src/shared/errutil/must"
	"giggler-golang/src/shared/singleton"
)

var getLogger = singleton.New(func() loggerInterface {
	switch must.GetEnv("LOGGER") {
	case "stdout":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "nop":
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	default:
		panic(must.ErrInvalidEnvValue)
	}
})

type loggerInterface interface {
	Debug(message string, vals ...any)
	Info(message string, vals ...any)
	Warn(message string, vals ...any)
	Error(message string, vals ...any)
}

// TODO: add a separate type for the key-value pairs

func Debug(message string, vals ...any) { getLogger().Debug(message, vals...) }

func Info(message string, vals ...any) { getLogger().Info(message, vals...) }

func Warn(message string, vals ...any) { getLogger().Warn(message, vals...) }

func Error(message string, vals ...any) { getLogger().Error(message, vals...) }
