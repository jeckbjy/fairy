package fairy

func Trace(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_TRACE,format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_DEBUG,format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_INFO,format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_WARN,format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_ERROR,format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Write(log.LEVEL_FATAL,format, args...)
}