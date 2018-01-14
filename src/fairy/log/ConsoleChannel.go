package log

type ConsoleChannel struct {
	BaseChannel
}

func (self *ConsoleChannel) Write(msg *Message) {
	// color
}