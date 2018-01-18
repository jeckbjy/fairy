package log

type Channel interface {
	Name() string
	Open()
	Close()
	Write(msg *Message)
	SetProperty(key string, val interface{})
	SetEnable(enable bool)
	GetEnable() bool
}
