package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/abc-valera/giggler-golang/src/components/env"
)

var loggerVar = initLogger()

// loggerInterface is used to provide a simpler interface for logging
type loggerInterface interface {
	Debug(message string, vals ...any)
	Info(message string, vals ...any)
	Warn(message string, vals ...any)
	Error(message string, vals ...any)
}

func initLogger() loggerInterface {
	switch env.Load("LOGGER") {
	case "slog_stdout":
		return slog.New(slog.NewTextHandler(os.Stdout, nil))
	case "nop":
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

func Debug(message string, vals ...any) { loggerVar.Debug(message, vals...) }

func Info(message string, vals ...any) { loggerVar.Info(message, vals...) }

func Warn(message string, vals ...any) { loggerVar.Warn(message, vals...) }

func Error(message string, vals ...any) { loggerVar.Error(message, vals...) }
