package log

const (
	LEVEL_ALL   = -1
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_OFF
	LEVEL_MAX
)

var gLevelName = []string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

type Message struct {
	Level    int
	File     string
	FileName string
	Line     int
	Timetamp int64
	Text     string
	Output   string
}
