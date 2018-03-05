package fairy

import (
	"fairy/util"
	"fmt"
	"reflect"
)

var gRegistry *Registry

func GetRegistry() *Registry {
	util.Once(gRegistry, func() {
		gRegistry = NewRegistry()
	})

	return gRegistry
}

func NewRegistry() *Registry {
	registry := &Registry{}
	registry.idMap = make(IdMap)
	registry.nameMap = make(NameMap)
	registry.typeMap = make(TypeMap)
	return registry
}

// 注册消息
func RegisterMessage(msg interface{}, args ...interface{}) {
	GetRegistry().Register(msg, args...)
}

// 消息元信息
type MsgInfo struct {
	Id   uint
	Name string
	Type reflect.Type
}

func (self *MsgInfo) New() interface{} {
	return reflect.New(self.Type).Interface()
}

type IdMap map[uint]*MsgInfo
type NameMap map[string]*MsgInfo
type TypeMap map[reflect.Type]*MsgInfo

// 消息注册：name->type,id->type,type->(name,id)
type Registry struct {
	idMap   IdMap
	nameMap NameMap
	typeMap TypeMap
}

func (self *Registry) Register(msg interface{}, args ...interface{}) error {
	if len(args) == 0 {
		return self.RegisterByName(msg, "")
	}

	if val, ok := args[0].(string); ok {
		return self.RegisterByName(msg, val)
	}

	id, err := util.ConvUint(args[0])
	if err != nil {
		return err
	}

	return self.RegisterById(msg, id)
}

func (self *Registry) RegisterByName(msg interface{}, msgName string) error {
	msgType := util.GetRealType(msg)

	if msgName == "" {
		msgName = msgType.Name()
	}

	if _, ok := self.typeMap[msgType]; ok {
		return fmt.Errorf("msg_type has registered![msg_name=%s]", msgName)
	}

	if _, ok := self.nameMap[msgName]; ok {
		return fmt.Errorf("msg_name has registered![msg_name=%s]", msgName)
	}

	info := &MsgInfo{Id: 0, Name: msgName, Type: msgType}

	self.typeMap[msgType] = info
	self.nameMap[msgName] = info

	return nil
}

func (self *Registry) RegisterById(msg interface{}, msgId uint) error {
	msgType := util.GetRealType(msg)
	msgName := msgType.Name()

	if msgId <= 0 {
		return fmt.Errorf("msgid must be greator than zero")
	}

	if _, ok := self.typeMap[msgType]; ok {
		return fmt.Errorf("msg_type has registered![msg_name=%s]", msgName)
	}

	if _, ok := self.nameMap[msgName]; ok {
		return fmt.Errorf("msg_name has registered![msg_name=%s]", msgName)
	}

	if _, ok := self.idMap[msgId]; ok {
		return fmt.Errorf("msg_id has registered![msg_name=%s, msg_id=%v]", msgName, msgId)
	}

	info := &MsgInfo{Id: msgId, Name: msgName, Type: msgType}
	self.typeMap[msgType] = info
	self.nameMap[msgName] = info
	self.idMap[msgId] = info

	return nil
}

func (self *Registry) Remove(msg interface{}) bool {
	msgType := util.GetRealType(msg)
	info, ok := self.typeMap[msgType]
	if ok {
		delete(self.typeMap, msgType)
		delete(self.nameMap, info.Name)
		if info.Id != 0 {
			delete(self.idMap, info.Id)
		}

		return true
	}

	return false
}

func (self *Registry) GetInfo(msg interface{}) (uint, string) {
	msgType := util.GetRealType(msg)
	info, ok := self.typeMap[msgType]
	if ok {
		return info.Id, info.Name
	}

	return 0, ""
}

func (self *Registry) GetName(msg interface{}) string {
	msgType := util.GetRealType(msg)
	info, ok := self.typeMap[msgType]
	if ok {
		return info.Name
	}

	return ""
}

func (self *Registry) GetId(msg interface{}) uint {
	msgType := util.GetRealType(msg)
	info, ok := self.typeMap[msgType]
	if ok {
		return info.Id
	}

	return 0
}

func (self *Registry) Create(id uint, name string) interface{} {
	if id == 0 {
		return self.CreateByName(name)
	} else {
		return self.CreateById(id)
	}
}

func (self *Registry) CreateByName(name string) interface{} {
	info, ok := self.nameMap[name]
	if ok {
		return info.New()
	}

	return nil
}

func (self *Registry) CreateById(id uint) interface{} {
	info, ok := self.idMap[id]
	if ok {
		return info.New()
	}

	return nil
}
