package log

import (
	"fmt"
	"strconv"
	"strings"
)

func NewConsoleChannel() *ConsoleChannel {
	channel := &ConsoleChannel{}
	channel.Init()
	channel.colorEnable = true
	channel.colorLevels[LEVEL_TRACE] = None
	channel.colorLevels[LEVEL_DEBUG] = None
	channel.colorLevels[LEVEL_INFO] = None
	channel.colorLevels[LEVEL_WARN] = Yellow
	channel.colorLevels[LEVEL_ERROR] = Red
	channel.colorLevels[LEVEL_FATAL] = Red
	return channel
}

type ConsoleChannel struct {
	BaseChannel
	colorEnable bool
	colorLevels [LEVEL_NUM]Color
}

func (*ConsoleChannel) Name() string {
	return "Console"
}

func (self *ConsoleChannel) Write(msg *Message) {
	// color
	output := self.GetOutput(msg)
	if self.colorEnable {
		Foreground(self.colorLevels[msg.Level])
		fmt.Print(output)
		Reset()
	} else {
		fmt.Print(output)
	}
}

func (self *ConsoleChannel) SetProperty(key string, val string) bool {
	if self.BaseChannel.SetProperty(key, val) {
		return true
	}

	if strings.HasPrefix(key, "color.") {
		key = key[6:]
		if key == "enable" {
			// color.enable = true
			enable, ok := strconv.ParseBool(val)
			if ok == nil {
				self.colorEnable = enable
				return true
			}
		} else {
			// color.debug = red
			level, lok := ParseLevel(key)
			color, cok := Parse(val)
			if lok && cok {
				self.colorLevels[level] = color
				return true
			}
		}

	}

	return false
}
