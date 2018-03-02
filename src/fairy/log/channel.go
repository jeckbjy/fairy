package log

type Channel interface {
	Name() string
	Open()
	Close()
	Write(msg *Message)
	SetProperty(key string, val string) bool
	GetConfig() *Config
	GetLogger() *Logger
	SetLogger(logger *Logger) 
}
