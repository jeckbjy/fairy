package log

func NewConsoleChannel() *ConsoleChannel {
	channel := &ConsoleChannel{}
	return channel
}

type ConsoleChannel struct {
	BaseChannel
}

func (self *ConsoleChannel) Name() string {
	return "Console"
}

func (self *ConsoleChannel) Write(msg *Message) {
	// color
}
