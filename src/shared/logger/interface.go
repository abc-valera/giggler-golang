package logger

// Interface is used to provide a simpler interface for logging
type Interface interface {
	Debug(message string, vals ...any)
	Info(message string, vals ...any)
	Warn(message string, vals ...any)
	Error(message string, vals ...any)
}
