package log

import "fmt"

// Message for logger channel
type Message struct {
	Level    int
	File     string
	FileName string
	Line     int
	Timetamp int64
	Text     string
	Output   string            // final text
	Data     map[string]string // 扩展数据
}

// Option 自定义Key-Value
type Option struct {
	Key string
	Val string
}

// WithOption example: WithOption("uid", "xxxxx")
func WithOption(key string, val interface{}) *Option {
	op := &Option{Key: key, Val: fmt.Sprintf("%+v", val)}
	return op
}
