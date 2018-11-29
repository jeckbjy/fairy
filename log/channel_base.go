package log

type BaseChannel struct {
	Config
	logger *Logger
}

func (*BaseChannel) Open() {

}

func (*BaseChannel) Close() {

}

func (bc *BaseChannel) SetProperty(key string, val string) bool {
	if bc.Config.SetConfig(key, val) {
		return true
	}

	return false
}

func (bc *BaseChannel) GetConfig() *Config {
	return &bc.Config
}

func (bc *BaseChannel) GetLogger() *Logger {
	return bc.logger
}

func (bc *BaseChannel) SetLogger(logger *Logger) {
	bc.logger = logger
}

func (bc *BaseChannel) GetOutput(msg *Message) string {
	if bc.Config.pattern != nil {
		return bc.Config.pattern.Format(msg)
	}

	return msg.Output
}
