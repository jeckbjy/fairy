package fairy

type FilterContext interface {
	AttrMap
	GetConn() Conn
	SetMessage(msg interface{})
	GetMessage() interface{}
	GetHandler() Handler
	SetHandler(handler Handler)
	GetError() error
	ThrowError(err error) FilterAction
	GetStopAction() FilterAction
	GetNextAction() FilterAction
	GetLastAction() FilterAction
	GetFirstAction() FilterAction
}
