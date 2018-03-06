package log

import (
	"fairy/util"
)

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

// func Traceu(uid uint64, format string) {

// }

func Catch() {
	if err := recover(); err != nil {
		Error("%+v", err)
		Error("%+v", util.GetStackTrace())
	}
}
