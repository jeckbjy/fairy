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
	return dispatcher
}

const (
	INVOKER_ID_MAX   = 65535 // 最大个数
	INVOKER_ID_LIMIT = 2048  // 这个值以下使用数组存储
)

type HandlerNameMap map[string]Handler

// Dispatcher 保存Handler的映射关系,以及middleware
type Dispatcher struct {
	nameMap     HandlerNameMap // 名字到Handler映射
	idArray     []Handler      // Handler数组,id是下标
	middlewares HandlerChain   // 中间件集合
}

func (self *Dispatcher) GetNameMap() HandlerNameMap {
	return self.nameMap
}

func (self *Dispatcher) GetIDArray() []Handler {
	return self.idArray
}

// Use 添加middleware
func (self *Dispatcher) Use(cb ...HandlerCB) {
	self.middlewares = append(self.middlewares, cb...)
}

// Middlewares 返回所有中间件
func (self *Dispatcher) Middlewares() HandlerChain {
	return self.middlewares
}

// Regsiter key:int,uint，string或者类，数字代表id查找，字符串或者类代表用名字查找
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

	if int(id) >= len(self.idArray) {
		targets := make([]Handler, id+1)
		copy(targets, self.idArray)
		self.idArray = targets
	}

	self.idArray[id] = handler

	return nil
}

func (self *Dispatcher) GetHandler(id uint, name string) Handler {
	if id > 0 {
		return self.GetHandlerById(id)
	}

	return self.GetHandlerByName(name)
}

// GetHandlerById 通过ID查询
func (self *Dispatcher) GetHandlerById(id uint) Handler {
	if id == 0 || id >= uint(len(self.idArray)) {
		return nil
	}

	return self.idArray[id]
}

// GetHandlerByName 通过名字查询
func (self *Dispatcher) GetHandlerByName(name string) Handler {
	return self.nameMap[name]
}
