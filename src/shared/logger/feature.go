package logger

func Debug(message string, vals ...any) { loggerVar.Debug(message, vals...) }

func Info(message string, vals ...any) { loggerVar.Info(message, vals...) }

func Warn(message string, vals ...any) { loggerVar.Warn(message, vals...) }

func Error(message string, vals ...any) { loggerVar.Error(message, vals...) }
