package base

import (
	"fairy"
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
	handler     fairy.Handler
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

func (self *BaseFilterContext) GetError() error {
	return self.err
}

func (self *BaseFilterContext) SetError(err error) {
	self.err = err
}

func (self *BaseFilterContext) ThrowError(err error) {
	self.err = err
	self.filterChain.HandleError(self.conn, err)
}

func (self *BaseFilterContext) SetHandler(handler fairy.Handler) {
	self.handler = handler
}

func (self *BaseFilterContext) GetHandler() fairy.Handler {
	return self.handler
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
