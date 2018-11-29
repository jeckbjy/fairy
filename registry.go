package fairy

import (
	"errors"
	"fmt"
	"reflect"
)

var gRegistry = NewRegistry()

// GetRegistry 获取全局Registry
func GetRegistry() *Registry {
	return gRegistry
}

// NewRegistry 创建Registry
func NewRegistry() *Registry {
	r := &Registry{}
	r.idMap = make(map[uint]*msgInfo)
	r.nameMap = make(map[string]*msgInfo)
	r.typeMap = make(map[reflect.Type]*msgInfo)
	return r
}

// RegisterMessage 注册消息
func RegisterMessage(msg interface{}, key interface{}) {
	GetRegistry().Register(msg, key)
}

// msgInfo 消息信息
type msgInfo struct {
	Id   uint
	Name string
	Type reflect.Type
}

// Registry 用于注册消息的id和name
type Registry struct {
	idMap   map[uint]*msgInfo
	nameMap map[string]*msgInfo
	typeMap map[reflect.Type]*msgInfo
}

// Register 注册消息元信息,id为字符串或者整数,nil则反射msg名字注册
func (r *Registry) Register(msg interface{}, key interface{}) error {
	msgid := uint(0)
	name := ""
	rtype := reflect.TypeOf(msg)
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	if key == nil {
		name = rtype.Name()
	} else {
		switch key.(type) {
		case string:
			name = key.(string)
		case uint:
			msgid = key.(uint)
		case int:
			msgid = uint(key.(int))
		case int32:
			msgid = uint(key.(int32))
		case uint32:
			msgid = uint(key.(uint32))
		case int16:
			msgid = uint(key.(int16))
		case uint16:
			msgid = uint(key.(uint16))
		default:
			return errors.New("register msg,bad type")
		}
	}

	if _, ok := r.typeMap[rtype]; ok {
		return fmt.Errorf("msg has registered,[msg=%+s]", rtype.Name())
	}

	if name == "" {
		name = rtype.Name()
	}

	info := &msgInfo{Id: msgid, Name: name, Type: rtype}
	r.typeMap[rtype] = info
	r.nameMap[name] = info

	if msgid != 0 {
		r.idMap[msgid] = info
	}

	return nil
}

// Remove 删除消息注册
func (r *Registry) Remove(msg interface{}) bool {
	mtype := reflect.TypeOf(msg)
	if mtype.Kind() == reflect.Ptr {
		mtype = mtype.Elem()
	}
	if info, ok := r.typeMap[mtype]; ok {
		delete(r.typeMap, mtype)
		delete(r.nameMap, info.Name)
		if info.Id != 0 {
			delete(r.idMap, info.Id)
		}

		return true
	}

	return false
}

// Create 通过Id或者名字创建
func (r *Registry) Create(msgid uint, name string) interface{} {
	var info *msgInfo
	if msgid != 0 {
		info = r.idMap[msgid]
	}

	if info == nil && name != "" {
		info = r.nameMap[name]
	}

	if info != nil {
		return reflect.New(info.Type).Interface()
	}

	return nil
}

// GetInfo 通过消息类型,查询信息
func (r *Registry) GetInfo(msg interface{}) (uint, string, bool) {
	mtype := reflect.TypeOf(msg)
	if mtype.Kind() == reflect.Ptr {
		mtype = mtype.Elem()
	}

	if info, ok := r.typeMap[mtype]; ok {
		return info.Id, info.Name, true
	}

	return 0, "", false
}
