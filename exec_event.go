package fairy

type Event interface {
	Process()
}

//////////////////////////////////////////////////
// FuncEvent
//////////////////////////////////////////////////

type Callback func()

func NewFuncEvent(cb Callback) *FuncEvent {
	ev := &FuncEvent{cb: cb}
	return ev
}

type FuncEvent struct {
	cb Callback
}

func (self *FuncEvent) Process() {
	self.cb()
}
