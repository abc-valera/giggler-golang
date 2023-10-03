package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/abc-valera/giggler-golang/src/components/env"
)

var loggerVar Interface

func init() {
	switch env.Load("LOGGER") {
	case "slog_stdout":
		loggerVar = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case "nop":
		loggerVar = slog.New(slog.NewTextHandler(io.Discard, nil))
	default:
		panic(env.ErrInvalidEnvValue)
	}
}
