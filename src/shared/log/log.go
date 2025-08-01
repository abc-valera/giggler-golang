package log

import (
	"io"
	"log/slog"
	"os"

	"giggler-golang/src/shared/env"
)

var loggerVar = func() loggerInterface {
	switch env.Load("LOGGER") {
	case "stdout":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "nop":
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	default:
		panic(env.ErrInvalidEnvValue)
	}
}()

type loggerInterface interface {
	Debug(message string, vals ...any)
	Info(message string, vals ...any)
	Warn(message string, vals ...any)
	Error(message string, vals ...any)
}

func Debug(message string, vals ...any) { loggerVar.Debug(message, vals...) }

func Info(message string, vals ...any) { loggerVar.Info(message, vals...) }

func Warn(message string, vals ...any) { loggerVar.Warn(message, vals...) }

func Error(message string, vals ...any) { loggerVar.Error(message, vals...) }
