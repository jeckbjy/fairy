package log

type BaseChannel struct {
	Config
	logger *Logger
}

func (self *BaseChannel) Open() {

}

func (self *BaseChannel) Close() {

}

func (self *BaseChannel) SetProperty(key string, val string) bool {
	if self.Config.SetConfig(key, val) {
		return true
	}

	return false
}

func (self *BaseChannel) GetConfig() *Config {
	return &self.Config
}

func (self *BaseChannel) GetLogger() *Logger {
	return self.logger
}

func (self *BaseChannel) SetLogger(logger *Logger) {
	self.logger = logger
}

func (self *BaseChannel) GetOutput(msg *Message) string {
	if self.Config.pattern != nil {
		return self.Config.pattern.Format(msg)
	}

	return msg.Output
}
