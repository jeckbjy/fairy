package log

import "github.com/jeckbjy/fairy/util"

// catch log
func Catch() {
	if err := recover(); err != nil {
		Error("%+v", err)
		Error("%+v", util.GetStackTrace())
	}
}

func Trace(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_TRACE, "", format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_DEBUG, "", format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_INFO, "", format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_WARN, "", format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_ERROR, "", format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Write(LEVEL_FATAL, "", format, args...)
}

///////////////////////////////////////////////////////////////
// with uid
///////////////////////////////////////////////////////////////
func Traceu(uid interface{}, format string, args ...interface{}) {
	GetLogger().Write(LEVEL_TRACE, util.ConvStr(uid), format, args...)
}

func Debugu(uid interface{}, format string, args ...interface{}) {
	GetLogger().Write(LEVEL_DEBUG, util.ConvStr(uid), format, args...)
}

func Infou(uid interface{}, format string, args ...interface{}) {
	GetLogger().Write(LEVEL_INFO, util.ConvStr(uid), format, args...)
}

func Warnu(uid interface{}, format string, args ...interface{}) {
	GetLogger().Write(LEVEL_WARN, util.ConvStr(uid), format, args...)
}

func Erroru(uid interface{}, format string, args ...interface{}) {
	GetLogger().Write(LEVEL_ERROR, util.ConvStr(uid), format, args...)
}

func Fatalu(uid interface{}, format string, args ...interface{}) {
	GetLogger().Write(LEVEL_FATAL, util.ConvStr(uid), format, args...)
}
