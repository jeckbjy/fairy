package base

import (
	"fairy"
	"net"
)

func NewContext(filterChain fairy.FilterChain, conn fairy.Connection) *BaseFilterContext {
	ctx := &BaseFilterContext{}
	ctx.conn = conn
	ctx.filterChain = filterChain
	return ctx
}

type BaseFilterContext struct {
	BaseAttrMap
	filterChain fairy.FilterChain
	conn        fairy.Connection
	message     interface{}
	err         error
}

func (self *BaseFilterContext) GetConnection() fairy.Connection {
	return self.conn
}

func (self *BaseFilterContext) SetMessage(msg interface{}) {
	self.message = msg
}

func (self *BaseFilterContext) GetMessage() interface{} {
	return self.message
}

func (self *BaseFilterContext) GetBuffer() *fairy.Buffer {
	return nil
}

func (self *BaseFilterContext) GetAddress() net.Addr {
	return nil
}

func (self *BaseFilterContext) GetError() error {
	return self.err
}

func (self *BaseFilterContext) SetError(err error) {
	self.err = err
}

func (self *BaseFilterContext) GetStopAction() fairy.FilterAction {
	return gStopAction
}

func (self *BaseFilterContext) GetNextAction() fairy.FilterAction {
	return gNextAction
}

func (self *BaseFilterContext) GetLastAction() fairy.FilterAction {
	return gLastAction
}

func (self *BaseFilterContext) GetFirstAction() fairy.FilterAction {
	return gFirstAction
}
