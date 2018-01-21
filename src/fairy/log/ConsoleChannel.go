package log

import (
	"fmt"
)

func NewConsoleChannel() *ConsoleChannel {
	channel := &ConsoleChannel{}
	channel.Init()
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
	output := self.GetOutput(msg)
	fmt.Print(output)
}
