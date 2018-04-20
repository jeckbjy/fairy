package log

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/jeckbjy/fairy/util"
)

var gLogger *Logger

// GetLogger represents get the golbal logger
func GetLogger() *Logger {
	if gLogger == nil {
		gLogger = NewLogger()
	}

	return gLogger
}

// NewLogger represents create new logger
func NewLogger() *Logger {
	logger := &Logger{}
	logger.Config.Init()
	logger.Config.SetFormat(DEFAULT_PATTERN)
	logger.cfgPath = "./logger.cfg"
	logger.sync = false
	logger.callerSkip = 2
	logger.stopped = true
	logger.messages = list.New()
	logger.mutex = &sync.Mutex{}
	logger.cond = sync.NewCond(logger.mutex)
	return logger
}

// 线程安全注意：默认异步执行，所有配置需要在Start前执行,只有Message队列是线程安全的
type Logger struct {
	Config
	cfgPath    string      // 配置文件路径，默认./logger.cfg
	sync       bool        // 是否同步
	callerSkip int         // 忽略堆栈
	stopped    bool        // 是否运行
	channels   []Channel   // 所有channel
	messages   *list.List  // 消息队列
	mutex      *sync.Mutex // mutex
	cond       *sync.Cond  // 用于线程同步
}

// SetCallerSkip represents set caller depth for skip
func (l *Logger) SetCallerSkip(skip int) {
	l.callerSkip = skip
}

// AddDefaultChannels add default channels
func (l *Logger) AddDefaultChannels() {
	if len(l.channels) == 0 {
		l.mutex.Lock()
		if len(l.channels) == 0 {
			l.AddChannel(NewConsoleChannel())
			l.AddChannel(NewFileChannel())
		}
		l.mutex.Unlock()
	}
}

// AddChannel add channel
func (l *Logger) AddChannel(channel Channel) {
	channel.Open()
	l.channels = append(l.channels, channel)
}

// DelChannel remove channel
func (l *Logger) DelChannel(name string) {
	for i, channel := range l.channels {
		if strings.EqualFold(channel.Name(), name) {
			channel.Close()
			l.channels = append(l.channels[:i], l.channels[i+1:]...)
			break
		}
	}
}

// GetChannel : find channel
func (l *Logger) GetChannel(name string) Channel {
	for _, channel := range l.channels {
		if strings.EqualFold(channel.Name(), name) {
			return channel
		}
	}

	return nil
}

/*
Load logger.cfg line by line
prefix # is comment
example:
# logger config
logger.level = debug
logger.format = [%y-%m-%d %H:%M:%S][%q][%U:%u][%t]

# console config
console.level = error

# file config
file.path = ./server.log
*/
func (l *Logger) Load() error {
	// 加载配置文件
	file, err := os.Open(l.cfgPath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 || text[0] == '#' {
			continue
		}

		index := strings.Index(text, "=")
		if index == -1 {
			continue
		}
		key := strings.TrimSpace(text[:index])
		val := strings.TrimSpace(text[index+1:])
		l.SetProperty(key, val)
	}

	file.Close()
	return nil
}

/*
example:
global: prefix . or logger
SetProperty(".level", "debug")
SetProterty("logger.level", "debug")
channel:
SetProperty("file.level", "debug")
SetProperty("file.path", "./server.log")
*/
func (l *Logger) SetProperty(key string, val string) error {
	index := strings.Index(key, ".")
	if index == -1 {
		return fmt.Errorf("no channel:key=%+v", key)
	}

	name := strings.ToLower(key[0:index])
	prop := strings.ToLower(key[index+1:])

	if index == 0 || name == "logger" {
		// logger global config
		if l.Config.SetConfig(prop, val) {
			return nil
		}

		switch prop {
		case "sync":
			l.sync, _ = strconv.ParseBool(val)
			return nil
		}

		return fmt.Errorf("not find prop[%+v] in logger", prop)
	}

	// set channel property
	l.AddDefaultChannels()
	channel := l.GetChannel(name)
	if channel == nil {
		return fmt.Errorf("not find channel:%+v", name)
	}

	if !channel.SetProperty(prop, val) {
		return fmt.Errorf("not find prop[%+v] in channel[%+v]", prop, name)
	}

	return nil
}

func (l *Logger) Write(level int, uid string, format string, args ...interface{}) {
	if !l.Enable || int(level) < l.Level {
		// not open
		return
	}

	_, file, line, ok := runtime.Caller(l.callerSkip)
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
	msg.Level = int(level)
	msg.Text = text
	msg.File = file
	msg.Line = line
	msg.Timetamp = util.Now()
	msg.FileName = fileName
	msg.Uid = uid

	if l.Config.pattern != nil {
		msg.Output = l.Config.pattern.Format(msg)
	}

	l.AddDefaultChannels()

	if l.sync {
		l.output(msg)
	} else {
		// lazy start
		l.Start()
		// sync thread
		l.mutex.Lock()
		l.messages.PushBack(msg)
		l.cond.Signal()
		l.mutex.Unlock()
	}
}

func (l *Logger) output(msg *Message) {
	for _, channel := range l.channels {
		cfg := channel.GetConfig()
		if cfg.Enable && msg.Level >= cfg.Level {
			channel.Write(msg)
		}
	}
}

// Run logger loop
func (l *Logger) Run() {
	for !l.stopped {
		messages := list.List{}
		l.mutex.Lock()
		for !l.stopped && l.messages.Len() == 0 {
			l.cond.Wait()
		}
		util.SwapList(&messages, l.messages)
		l.mutex.Unlock()
		// process
		for iter := messages.Front(); iter != nil; iter = iter.Next() {
			msg := iter.Value.(*Message)
			l.output(msg)
		}
	}

	// close all channel
	for _, channel := range l.channels {
		channel.Close()
	}
}

// Start logger
func (l *Logger) Start() {
	if l.stopped {
		l.stopped = false
		go l.Run()
	}
}

// Stop logger
func (l *Logger) Stop() {
	if !l.stopped {
		l.stopped = true
		l.mutex.Lock()
		l.cond.Signal()
		l.mutex.Unlock()
	}
}
