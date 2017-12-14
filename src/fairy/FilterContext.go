package fairy

import "net"

type FilterContext interface {
	AttrMap
	GetConnection() Connection
	// SetData(data interface{})
	// GetData() interface{}
	SetMessage(msg interface{})
	GetMessage() interface{}
	GetAddress() net.Addr
	GetStopAction() FilterAction
	GetNextAction() FilterAction
	GetLastAction() FilterAction
	GetFirstAction() FilterAction
}
