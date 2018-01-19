package log

type Channel interface {
	Name() string
	Open()
	Close()
	Write(msg *Message)
	SetProperty(key string, val interface{}) bool
	GetConfig() *Config
}
