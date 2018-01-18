package log

const (
	LEVEL_ALL = iota
	LEVEL_TRACE
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_OFF
	LEVEL_MAX
)

type Message struct {
	Level    int
	Info     string
	File     string
	Line     int
	Timetamp int64
}
