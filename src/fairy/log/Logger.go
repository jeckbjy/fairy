package log

var gLogger *Logger

func GetLogger() *Logger {
	if gLogger == nil {
		gLogger = NewLogger()
		// set default
	}

	return gLogger
}

func NewLogger() *Logger {
	logger := &Logger{}
	return logger
}

type Logger struct {
	channels []Channel
}

func (self *Logger) AddChannel(channel Channel) {
	self.channels = append(self.channels, channel)
}

func (self *Logger) DelChannel(name string) {
	// del by name
}

func (self *Logger) SetProperty(key string, val interface{}) {

}

func (self *Logger) Write(level int, format string, args ...interface{}) {

}

func (self *Logger) Run() {

}

func (self *Logger) Start() {

}

func (self *Logger) Stop() {

}
