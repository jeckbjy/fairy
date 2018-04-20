package log

// Message for logger channel
type Message struct {
	Level    int
	File     string
	FileName string
	Line     int
	Timetamp int64
	Text     string
	Uid      string // TODO:support uid
	Output   string // 最终结果
}
