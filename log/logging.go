package log

import (
	"runtime"
)

// Catch log
func Catch() {
	if err := recover(); err != nil {
		buf := make([]byte, 1<<15)
		stacklen := runtime.Stack(buf, false)
		stack := string(buf[:stacklen])
		Error("%+v", err)
		Error("%+v", stack)
	}
}

func Trace(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_TRACE, format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_FATAL, format, args...)
}
