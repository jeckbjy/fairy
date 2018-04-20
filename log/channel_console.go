package log

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jeckbjy/fairy/util/terminal"
)

func NewConsoleChannel() *ConsoleChannel {
	channel := &ConsoleChannel{}
	channel.Init()
	channel.colorEnable = true
	channel.colorLevels[LEVEL_TRACE] = terminal.None
	channel.colorLevels[LEVEL_DEBUG] = terminal.None
	channel.colorLevels[LEVEL_INFO] = terminal.None
	channel.colorLevels[LEVEL_WARN] = terminal.Yellow
	channel.colorLevels[LEVEL_ERROR] = terminal.Red
	channel.colorLevels[LEVEL_FATAL] = terminal.Red
	return channel
}

type ConsoleChannel struct {
	BaseChannel
	colorEnable bool
	colorLevels [LEVEL_NUM]terminal.Color
}

func (self *ConsoleChannel) Name() string {
	return "Console"
}

func (self *ConsoleChannel) Write(msg *Message) {
	// color
	output := self.GetOutput(msg)
	if self.colorEnable {
		terminal.Foreground(self.colorLevels[msg.Level])
		fmt.Print(output)
		terminal.Reset()
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
			color, cok := terminal.Parse(val)
			if lok && cok {
				self.colorLevels[level] = color
				return true
			}
		}

	}

	return false
}
