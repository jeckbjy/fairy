package fairy

import (
	"fmt"
	"reflect"

	"github.com/jeckbjy/fairy/util"
)

var gDispatcher *Dispatcher

func GetDispatcher() *Dispatcher {
	util.Once(gDispatcher, func() {
		gDispatcher = NewDispatcher()
	})

	return gDispatcher
}

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{}
	dispatcher.nameMap = make(HandlerNameMap)
	dispatcher.idMap = make(HandlerIdMap)
	return dispatcher
}

type HandlerNameMap map[string]Handler
type HandlerIdMap map[uint]Handler

const (
	INVOKER_ID_MAX   = 65535 // 最大个数
	INVOKER_ID_LIMIT = 2048  // 这个值以下使用数组存储
)

// Dispatcher 保存Handler的映射关系
type Dispatcher struct {
	nameMap  HandlerNameMap
	idMap    HandlerIdMap
	idArray  []Handler
	uncaught Handler
}

func (self *Dispatcher) GetNameMap() HandlerNameMap {
	return self.nameMap
}

func (self *Dispatcher) GetIDMap() HandlerIdMap {
	return self.idMap
}

func (self *Dispatcher) GetIDArray() []Handler {
	return self.idArray
}

// key:int,uint，string或者类，数字代表id查找，字符串或者类代表用名字查找
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

	return nil, false
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
