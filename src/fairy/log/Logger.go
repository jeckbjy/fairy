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
	logger.callerSkip = 1
	return logger
}

type Logger struct {
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

func (self *Logger) SetProperty(key string, val interface{}) {

}

func (self *Logger) Write(level int, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(self.callerSkip)
	if !ok {
		return
	}

	info := fmt.Sprintf(format, args...)

	// fmt.Printf("%+v,%+v, %+v\n", file, line, info)

	msg := &Message{}
	msg.Level = level
	msg.Info = info
	msg.File = file
	msg.Line = line
	msg.Timetamp = util.Now()

	self.mutex.Lock()
	self.messages.PushBack(msg)
	self.cond.Signal()
	self.mutex.Unlock()
}

func (self *Logger) Run() {
	for {
		self.mutex.Lock()
		for !self.stopped && self.messages.Len() == 0 {
			self.cond.Wait()
		}
		self.mutex.Unlock()
		// process
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
