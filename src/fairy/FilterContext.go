package fairy

type FilterContext interface {
	AttrMap
	GetConnection() Connection
	SetMessage(msg interface{})
	GetMessage() interface{}
	GetHandler() Handler
	SetHandler(handler Handler)
	GetStopAction() FilterAction
	GetNextAction() FilterAction
	GetLastAction() FilterAction
	GetFirstAction() FilterAction
}
