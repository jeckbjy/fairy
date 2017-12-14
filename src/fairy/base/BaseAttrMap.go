package base

import (
	. "fairy"
)

type BaseAttrMap struct {
	attrs []interface{}
}

func (self *BaseAttrMap) HasAttr(key *AttrKey) bool {
	index := key.GetIndex()
	if index < len(self.attrs) && self.attrs[index] != nil {
		return true
	}

	return false
}

func (self *BaseAttrMap) SetAttr(key *AttrKey, val interface{}) {
	index := key.GetIndex()
	if index >= len(self.attrs) {
		count := index - len(self.attrs) + 1
		for i := 0; i < count; i++ {
			self.attrs = append(self.attrs, nil)
		}
	}

	self.attrs[index] = val
}

func (self *BaseAttrMap) GetAttr(key *AttrKey) interface{} {
	index := key.GetIndex()
	if index < len(self.attrs) {
		return self.attrs[index]
	}

	return nil
}

func (self *BaseAttrMap) GetAttrEx(key *AttrKey, defVal interface{}) interface{} {
	index := key.GetIndex()
	if index < len(self.attrs) {
		return self.attrs[index]
	}

	return defVal
}
