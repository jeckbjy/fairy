package log

import (
	"container/list"
	"fairy/util"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var gLogger *Logger

func GetLogger() *Logger {
	if gLogger == nil {
		gLogger = NewLogger()
		// set default
		gLogger.AddChannel(NewConsoleChannel())
		gLogger.AddChannel(NewFileChannel())
	}

	return gLogger
}

func NewLogger() *Logger {
	logger := &Logger{}
	logger.Config.Init()
	logger.Config.SetFormat(DEFAULT_PATTERN)
	logger.sync = false
	logger.messages = list.New()
	logger.mutex = &sync.Mutex{}
	logger.cond = sync.NewCond(logger.mutex)
	logger.stopped = true
	logger.callerSkip = 2
	return logger
}

type Logger struct {
	Config
	sync       bool
	channels   []Channel
	messages   *list.List
	mutex      *sync.Mutex
	cond       *sync.Cond
	stopped    bool
	callerSkip int
}

func (self *Logger) SetCallerSkip(skip int) {
	self.callerSkip = skip
}

func (self *Logger) AddChannel(channel Channel) {
	channel.Open()
	self.channels = append(self.channels, channel)
}

func (self *Logger) DelChannel(name string) {
	for i, channel := range self.channels {
		if strings.EqualFold(channel.Name(), name) {
			channel.Close()
			self.channels = append(self.channels[:i], self.channels[i+1:]...)
			break
		}
	}
}

func (self *Logger) GetChannel(name string) Channel {
	for _, channel := range self.channels {
		if strings.EqualFold(channel.Name(), name) {
			return channel
		}
	}

	return nil
}

/*
example:
global: prefix .
SetProperty(".level", "debug")
channel:
SetProperty("file.level", "debug")
SetProperty("file.path", "./server.log")
*/
func (self *Logger) SetProperty(key string, val string) error {
	index := strings.Index(key, ".")
	if index == -1 {
		return fmt.Errorf("no channel:key=%+v", key)
	}

	prop := strings.ToLower(key[index+1:])

	if index == 0 {
		// logger global config
		if self.Config.SetConfig(prop, val) {
			return nil
		}

		switch prop {
		case "sync":
			self.sync, _ = strconv.ParseBool(val)
			return nil
		}

		return fmt.Errorf("not find prop[%+v] in logger", prop)
	} else {
		name := key[0:index]
		channel := self.GetChannel(name)
		if channel == nil {
			return fmt.Errorf("not find channel:%+v", name)
		}

		if !channel.SetProperty(prop, val) {
			return fmt.Errorf("not find prop[%+v] in channel[%+v]", prop, name)
		}

		return nil
	}
}

func (self *Logger) Write(level int, format string, args ...interface{}) {
	if !self.Enable || level < self.Level {
		// not open
		return
	}

	_, file, line, ok := runtime.Caller(self.callerSkip)
	if !ok {
		return
	}

	var fileName string
	text := fmt.Sprintf(format, args...)
	index := strings.LastIndex(file, "/")
	if index != -1 {
		fileName = file[index+1:]
	}

	msg := &Message{}
	msg.Level = level
	msg.Text = text
	msg.File = file
	msg.Line = line
	msg.Timetamp = util.Now()
	msg.FileName = fileName

	if self.Config.pattern != nil {
		msg.Output = self.Config.pattern.Format(msg)
	}

	if self.sync {
		self.output(msg)
	} else {
		// lazy start
		self.Start()
		// sync thread
		self.mutex.Lock()
		self.messages.PushBack(msg)
		self.cond.Signal()
		self.mutex.Unlock()
	}
}

func (self *Logger) output(msg *Message) {
	for _, channel := range self.channels {
		cfg := channel.GetConfig()
		if cfg.Enable && msg.Level >= cfg.Level {
			channel.Write(msg)
		}
	}
}

func (self *Logger) Run() {
	for !self.stopped {
		messages := list.List{}
		self.mutex.Lock()
		for !self.stopped && self.messages.Len() == 0 {
			self.cond.Wait()
		}
		util.SwapList(&messages, self.messages)
		self.mutex.Unlock()
		// process
		for iter := messages.Front(); iter != nil; iter = iter.Next() {
			msg := iter.Value.(*Message)
			self.output(msg)
		}
	}

	// close all
	for _, channel := range self.channels {
		channel.Close()
	}
}

func (self *Logger) Start() {
	if self.stopped {
		self.stopped = false
		go self.Run()
	}
}

func (self *Logger) Stop() {
	if !self.stopped {
		self.stopped = true
		self.mutex.Lock()
		self.cond.Signal()
		self.mutex.Unlock()
	}
}
