package base

import (
	"fairy"
)

func NewContext(filterChain fairy.FilterChain, conn fairy.Connection) *FilterContext {
	ctx := &FilterContext{}
	ctx.conn = conn
	ctx.filterChain = filterChain
	return ctx
}

type FilterContext struct {
	AttrMap
	filterChain fairy.FilterChain
	conn        fairy.Connection
	message     interface{}
	handler     fairy.Handler
	err         error
}

func (self *FilterContext) GetConnection() fairy.Connection {
	return self.conn
}

func (self *FilterContext) SetMessage(msg interface{}) {
	self.message = msg
}

func (self *FilterContext) GetMessage() interface{} {
	return self.message
}

func (self *FilterContext) GetBuffer() *fairy.Buffer {
	return nil
}

func (self *FilterContext) GetError() error {
	return self.err
}

func (self *FilterContext) SetError(err error) {
	self.err = err
}

func (self *FilterContext) ThrowError(err error) {
	self.err = err
	self.filterChain.HandleError(self.conn, err)
}

func (self *FilterContext) SetHandler(handler fairy.Handler) {
	self.handler = handler
}

func (self *FilterContext) GetHandler() fairy.Handler {
	return self.handler
}

func (self *FilterContext) GetStopAction() fairy.FilterAction {
	return gStopAction
}

func (self *FilterContext) GetNextAction() fairy.FilterAction {
	return gNextAction
}

func (self *FilterContext) GetLastAction() fairy.FilterAction {
	return gLastAction
}

func (self *FilterContext) GetFirstAction() fairy.FilterAction {
	return gFirstAction
}
