package log

import "strings"

const (
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_OFF
	LEVEL_NUM = LEVEL_FATAL + 1
	LEVEL_ALL = -1
)

var gLevelName = []string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal", "Off"}

func ParseLevel(str string) (int, bool) {
	for i, name := range gLevelName {
		if strings.EqualFold(str, name) {
			return i, true
		}
	}

	return LEVEL_OFF, false
}
