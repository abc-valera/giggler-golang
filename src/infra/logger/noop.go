package logger

import (
	"io"
	"log/slog"
)

func InitNoop() Interface {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
