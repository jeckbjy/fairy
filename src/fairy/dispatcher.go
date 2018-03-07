package fairy

import (
	"fairy/log"
	"fairy/util"
	"fmt"
	"reflect"
)

var gErrorHandler Handler
var gDispatcher *Dispatcher

func GetDispatcher() *Dispatcher {
	util.Once(gDispatcher, func() {
		gDispatcher = NewDispatcher()
	})

	return gDispatcher
}

func error_cb(conn Conn, pkt Packet) {
	log.Error("cannot find handler:name=%+v,id=%+v,rpcid=%+v",
		pkt.GetName(),
		pkt.GetId(),
		pkt.GetRpcId())
}

func NewDispatcher() *Dispatcher {
	if gErrorHandler == nil {
		gErrorHandler = &HandlerHolder{cb: error_cb}
	}

	dispatcher := &Dispatcher{}
	dispatcher.nameMap = make(HandlerNameMap)
	dispatcher.idMap = make(HandlerIdMap)
	dispatcher.errhandler = gErrorHandler
	return dispatcher
}

//////////////////////////////////////////////////////////
// 注册回调函数
//////////////////////////////////////////////////////////
func RegisterHandler(key interface{}, cb HandlerCB) {
	RegisterHandlerEx(key, cb, 0)
}

func RegisterHandlerEx(key interface{}, cb HandlerCB, queueId int) {
	GetDispatcher().Regsiter(key, &HandlerHolder{cb: cb, queueId: queueId})
}

func RegisterUncaughtHandler(cb HandlerCB) {
	GetDispatcher().SetUncaughtHandler(&HandlerHolder{cb: cb, queueId: 0})
}

//////////////////////////////////////////////////////////
// Handler
//////////////////////////////////////////////////////////
type Handler interface {
	GetQueueId() int // 暗示在哪个线程执行
	Invoke(conn Conn, packet Packet)
}

type HandlerCB func(Conn, Packet)
type HandlerHolder struct {
	cb      HandlerCB
	queueId int
}

func (self *HandlerHolder) GetQueueId() int {
	return self.queueId
}

func (self *HandlerHolder) Invoke(conn Conn, packet Packet) {
	defer log.Catch()
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
	nameMap    HandlerNameMap
	idMap      HandlerIdMap
	idArray    []Handler
	uncaught   Handler
	errhandler Handler
}

/**
 * key:int,uint，string或者类，数字代表id查找，字符串或者类代表用名字查找
 */
func (self *Dispatcher) Regsiter(key interface{}, handler Handler) {
	switch key.(type) {
	case int:
		self.RegistryById(uint(key.(int)), handler)
	case uint:
		self.RegistryById(key.(uint), handler)
	case string:
		self.RegisterByName(key.(string), handler)
	default:
		// must be struct!!!
		// example:Register(&LoginReq{}, handler) or Register(LoginReq{}, handler)
		rtype := util.GetRealType(key)
		if rtype.Kind() == reflect.Struct {
			self.RegisterByName(rtype.Name(), handler)
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

func (self *Dispatcher) GetFinalHandler(id uint, name string) (Handler, bool) {
	h := self.GetHandler(id, name)
	if h != nil {
		return h, true
	}

	if self.uncaught != nil {
		return self.uncaught, false
	}

	return self.errhandler, false
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

func (self *Dispatcher) GetErrorHandler() Handler {
	return self.errhandler
}

func (self *Dispatcher) SetErrorHandler(h Handler) {
	self.errhandler = h
}
