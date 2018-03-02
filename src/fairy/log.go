package fairy

import (
	"fairy/log"
	"fairy/util"
)

func GetLogger() *log.Logger {
	logger := log.GetLogger()
	logger.SetCallerSkip(2)
	return logger
}

func Trace(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_TRACE, format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_FATAL, format, args...)
}

func Catch() {
	if err := recover(); err != nil {
		Error("%+v", err)
		Error("%+v", util.GetStackTrace())
	}
}
