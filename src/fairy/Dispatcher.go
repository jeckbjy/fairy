package fairy

import (
	"fairy/util"
	"fmt"
)

type Invoker interface {
	Invoke(conn Connection, packet Packet)
}

type NameInvokerMap map[string]Invoker
type IdInvokeMap map[int]Invoker

const (
	INVOKER_ID_MAX   = 65535 // 最大个数
	INVOKER_ID_LIMIT = 2048  // 这个值以下使用数组存储
)

type Dispatcher struct {
	nameInvokerMap NameInvokerMap
	idInvokerMap   IdInvokeMap
	idInvokerArray []Invoker
}

func (self *Dispatcher) Regsiter(key interface{}, invoker Invoker) {
	switch key.(type) {
	case string:
		self.RegisterByName(key.(string), invoker)
	case int:
		self.RegistryById(key.(int), invoker)
	case uint:
		self.RegistryById(int(key.(uint)), invoker)
	default:
		panic(fmt.Sprintf("register invoker fail,bad key type!key=%+v", key))
	}
}

func (self *Dispatcher) RegisterByName(name string, invoker Invoker) error {
	if _, ok := self.nameInvokerMap[name]; ok {
		return util.NewError("register invoker is duplicate! name=%s", name)
	}

	self.nameInvokerMap[name] = invoker
	return nil
}

func (self *Dispatcher) RegistryById(id int, invoker Invoker) error {
	if self.GetInvokerById(id) != nil {
		return util.NewError("invoker id has registered!id=%v", id)
	}

	if id < 0 || id > INVOKER_ID_MAX {
		return util.NewError("invoker id overflow!id=%v", id)
	}

	if id < INVOKER_ID_LIMIT {
		if id >= len(self.idInvokerArray) {
			// resize
			count := id - len(self.idInvokerArray) + 1
			self.idInvokerArray = append(self.idInvokerArray, make([]Invoker, count)...)
		}
		self.idInvokerArray[id] = invoker
	} else {
		self.idInvokerMap[id] = invoker
	}

	return nil
}

func (self *Dispatcher) GetInvokerById(id int) Invoker {
	if id < 0 || id > INVOKER_ID_MAX {
		return nil
	}

	if id < len(self.idInvokerArray) {
		return self.idInvokerArray[id]
	} else if id > INVOKER_ID_LIMIT {
		return self.idInvokerMap[id]
	}

	return nil
}
