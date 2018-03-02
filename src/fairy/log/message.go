package log

// Message for logger channel
type Message struct {
	Level    int
	File     string
	FileName string
	Line     int
	Timetamp int64
	Text     string
	Output   string
}
