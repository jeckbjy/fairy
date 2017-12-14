package fairy

import (
	"errors"
	"fmt"
	"reflect"
)

var gRegistry *Registry

func GetGlobalRegistry() *Registry {
	if gRegistry == nil {
		gRegistry = NewRegistry()
	}

	return gRegistry
}

func NewRegistry() *Registry {
	registry := &Registry{}
	registry.idMap = make(IdMap)
	registry.nameMap = make(NameMap)
	registry.typeMap = make(TypeMap)
	return registry
}

// 消息元信息
type MsgInfo struct {
	Id   uint
	Name string
	Type reflect.Type
}

type IdMap map[uint]*MsgInfo
type NameMap map[string]*MsgInfo
type TypeMap map[reflect.Type]*MsgInfo

// 消息分发
// type Dispatcher struct {
// }

// 消息注册：name->type,id->type,type->(name,id)
type Registry struct {
	idMap   IdMap
	nameMap NameMap
	typeMap TypeMap
}

func (self *Registry) Register(msg interface{}) error {
	msg_type := reflect.TypeOf(msg)
	msg_name := msg_type.Name()
	if _, ok := self.typeMap[msg_type]; ok {
		return errors.New(fmt.Sprintf("msg_type has registered![msg_name=%s]", msg_type.Name()))
	}

	if _, ok := self.nameMap[msg_name]; ok {
		return errors.New(fmt.Sprintf("msg_name has registered![msg_name=%s]", msg_type.Name()))
	}

	info := &MsgInfo{Id: 0, Name: msg_type.Name(), Type: msg_type}

	self.typeMap[msg_type] = info
	self.nameMap[msg_name] = info

	return nil
}

func (self *Registry) RegisterId(msg interface{}, msg_id uint) error {
	msg_type := reflect.TypeOf(msg)
	msg_name := msg_type.Name()

	if _, ok := self.typeMap[msg_type]; ok {
		return errors.New(fmt.Sprintf("msg_type has registered![msg_name=%s]", msg_name))
	}

	if _, ok := self.nameMap[msg_name]; ok {
		return errors.New(fmt.Sprintf("msg_name has registered![msg_name=%s]", msg_name))
	}

	if _, ok := self.idMap[msg_id]; ok {
		return errors.New(fmt.Sprintf("msg_id has registered![msg_name=%s, msg_id=%v]", msg_name, msg_id))
	}

	info := &MsgInfo{Id: msg_id, Name: msg_type.Name(), Type: msg_type}
	self.typeMap[msg_type] = info
	self.nameMap[msg_name] = info
	self.idMap[msg_id] = info

	return nil
}

func (self *Registry) Remove(msg interface{}) bool {
	msg_type := reflect.TypeOf(msg)
	info, ok := self.typeMap[msg_type]
	if ok {
		delete(self.typeMap, msg_type)
		delete(self.nameMap, info.Name)
		if info.Id != 0 {
			delete(self.idMap, info.Id)
		}

		return true
	} else {
		return false
	}
}

func (self *Registry) GetName(msg interface{}) string {
	msg_type := reflect.TypeOf(msg)
	info, ok := self.typeMap[msg_type]
	if ok {
		return info.Name
	} else {
		return ""
	}
}

func (self *Registry) GetId(msg interface{}) uint {
	msg_type := reflect.TypeOf(msg)
	info, ok := self.typeMap[msg_type]
	if ok {
		return info.Id
	} else {
		return 0
	}
}

func (self *Registry) CreateByName(name string) interface{} {
	info, ok := self.nameMap[name]
	if ok {
		return reflect.New(info.Type)
	} else {
		return nil
	}
}

func (self *Registry) CreateById(id uint) interface{} {
	info, ok := self.idMap[id]
	if ok {
		return reflect.New(info.Type)
	} else {
		return nil
	}
}
