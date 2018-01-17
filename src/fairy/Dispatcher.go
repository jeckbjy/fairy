package fairy

import (
	"fmt"
	"reflect"
)

var gDispatcher *Dispatcher

func GetDispatcher() *Dispatcher {
	if gDispatcher == nil {
		gDispatcher = NewDispatcher()
	}

	return gDispatcher
}

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{}
	dispatcher.nameMap = make(HandlerNameMap)
	dispatcher.idMap = make(HandlerIdMap)
	return dispatcher
}

//////////////////////////////////////////////////////////
// 注册回调函数
//////////////////////////////////////////////////////////
func RegisterHandler(key interface{}, cb HandlerCallback) {
	RegisterHandlerEx(key, cb, 0)
}

func RegisterHandlerEx(key interface{}, cb HandlerCallback, queueId int) {
	GetDispatcher().Regsiter(key, &HandlerHolder{cb: cb, queueId: queueId})
}

func RegisterUncaughtHandler(cb HandlerCallback) {
	GetDispatcher().SetUncaughtHandler(&HandlerHolder{cb: cb, queueId: 0})
}

//////////////////////////////////////////////////////////
// Handler
//////////////////////////////////////////////////////////
type Handler interface {
	GetQueueId() int // 暗示在哪个线程执行
	Invoke(conn Connection, packet Packet)
}

type HandlerCallback func(Connection, Packet)
type HandlerHolder struct {
	cb      HandlerCallback
	queueId int
}

func (self *HandlerHolder) GetQueueId() int {
	return self.queueId
}

func (self *HandlerHolder) Invoke(conn Connection, packet Packet) {
	self.cb(conn, packet)
}

//////////////////////////////////////////////////////////
// Dispatcher
//////////////////////////////////////////////////////////
type HandlerNameMap map[string]Handler
type HandlerIdMap map[uint]Handler

const (
	INVOKER_ID_MAX   = 65535 // 最大个数
	INVOKER_ID_LIMIT = 2048  // 这个值以下使用数组存储
)

type Dispatcher struct {
	nameMap  HandlerNameMap
	idMap    HandlerIdMap
	idArray  []Handler
	uncaught Handler
}

func (self *Dispatcher) Regsiter(key interface{}, handler Handler) {
	switch key.(type) {
	case string:
		self.RegisterByName(key.(string), handler)
	case int:
		self.RegistryById(key.(uint), handler)
	case uint:
		self.RegistryById(key.(uint), handler)
	default:
		// get name
		rtype := reflect.TypeOf(key)
		if rtype.Elem().Kind() == reflect.Struct {
			self.RegisterByName(rtype.Elem().Name(), handler)
		} else {
			panic(fmt.Sprintf("register handler fail,bad key type!key=%+v", key))
		}
	}
}

func (self *Dispatcher) RegisterByName(name string, handler Handler) error {
	if _, ok := self.nameMap[name]; ok {
		return fmt.Errorf("register invoker is duplicate! name=%s", name)
	}

	self.nameMap[name] = handler
	return nil
}

func (self *Dispatcher) RegistryById(id uint, handler Handler) error {
	if self.GetHandlerById(id) != nil {
		return fmt.Errorf("invoker id has registered!id=%v", id)
	}

	if id <= 0 || id > INVOKER_ID_MAX {
		return fmt.Errorf("invoker id overflow!id=%v", id)
	}

	if id < INVOKER_ID_LIMIT {
		idLen := uint(len(self.idArray))
		if id >= idLen {
			// resize
			count := id - idLen + 1
			self.idArray = append(self.idArray, make([]Handler, count)...)
		}
		self.idArray[id] = handler
	} else {
		self.idMap[id] = handler
	}

	return nil
}

func (self *Dispatcher) GetHandler(id uint, name string) Handler {
	if id > 0 {
		return self.GetHandlerById(id)
	}

	return self.GetHandlerByName(name)
}

func (self *Dispatcher) GetHandlerById(id uint) Handler {
	if id == 0 || id > INVOKER_ID_MAX {
		return nil
	}

	if id < uint(len(self.idArray)) {
		return self.idArray[id]
	} else if id > INVOKER_ID_LIMIT {
		return self.idMap[id]
	}

	return nil
}

func (self *Dispatcher) GetHandlerByName(name string) Handler {
	return self.nameMap[name]
}

func (self *Dispatcher) SetUncaughtHandler(handler Handler) {
	self.uncaught = handler
}

func (self *Dispatcher) GetUncaughtHandler() Handler {
	return self.uncaught
}
