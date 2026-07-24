package logger

// TODO: add a separate type for the key-value pairs

type Interface interface {
	Debug(message string, vals ...any)
	Info(message string, vals ...any)
	Warn(message string, vals ...any)
	Error(message string, vals ...any)
}
