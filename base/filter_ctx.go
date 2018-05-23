package base

import "github.com/jeckbjy/fairy"

func NewContext(filterChain fairy.FilterChain, conn fairy.Conn) *FilterContext {
	ctx := &FilterContext{}
	ctx.conn = conn
	ctx.filters = filterChain
	return ctx
}

type FilterContext struct {
	AttrMap
	filters fairy.FilterChain
	conn    fairy.Conn
	message interface{}
	err     error
}

func (self *FilterContext) GetConn() fairy.Conn {
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

func (self *FilterContext) ThrowError(err error) fairy.FilterAction {
	self.err = err
	self.filters.HandleError(self.conn, err)
	return self.GetStopAction()
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
