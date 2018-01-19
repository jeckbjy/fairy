package log

import (
	"container/list"
	"fairy/util"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var gLogger *Logger

func GetLogger() *Logger {
	if gLogger == nil {
		gLogger = NewLogger()
		gLogger.Config.Init()
		gLogger.SetFormat(DEFAULT_PATTERN)
		// set default
		gLogger.AddChannel(NewConsoleChannel())
		gLogger.AddChannel(NewFileChannel())
	}

	return gLogger
}

func NewLogger() *Logger {
	logger := &Logger{}
	logger.messages = list.New()
	logger.mutex = &sync.Mutex{}
	logger.cond = sync.NewCond(logger.mutex)
	logger.stopped = true
	logger.callerSkip = 2
	return logger
}

type Logger struct {
	Config
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
	self.channels = append(self.channels, channel)
}

func (self *Logger) DelChannel(name string) {
	for i, channel := range self.channels {
		if strings.EqualFold(channel.Name(), name) {
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

func (self *Logger) SetProperty(key string, val interface{}) {
	for _, channel := range self.channels {
		if channel.SetProperty(key, val) {
			break
		}
	}
}

func (self *Logger) SetChannelProperty(name string, key string, val interface{}) {
	if name == "" {
		self.SetProperty(key, val)
	}

	channel := self.GetChannel(name)
	if channel != nil {
		channel.SetProperty(key, val)
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

	// fmt.Printf("%+v,%+v, %+v, %+v\n", file, line, info, fileName)

	msg := &Message{}
	msg.Level = level
	msg.Text = text
	msg.File = file
	msg.Line = line
	msg.Timetamp = util.Now()
	msg.FileName = fileName

	tt := self.pattern.Format(msg)
	fmt.Printf("aaaa :%s", tt)

	// self.mutex.Lock()
	// self.messages.PushBack(msg)
	// self.cond.Signal()
	// self.mutex.Unlock()
}

func (self *Logger) Run() {
	for {
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
			for _, channel := range self.channels {
				cfg := channel.GetConfig()
				if cfg.Enable && msg.Level >= cfg.Level {
					channel.Write(msg)
				}
			}
		}
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
