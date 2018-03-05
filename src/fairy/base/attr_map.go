package base

import "fairy"

type AttrMap struct {
	attrs []interface{}
}

func (self *AttrMap) HasAttr(key *fairy.AttrKey) bool {
	index := key.GetIndex()
	if index < len(self.attrs) && self.attrs[index] != nil {
		return true
	}

	return false
}

func (self *AttrMap) SetAttr(key *fairy.AttrKey, val interface{}) {
	index := key.GetIndex()
	if index >= len(self.attrs) {
		count := index - len(self.attrs) + 1
		for i := 0; i < count; i++ {
			self.attrs = append(self.attrs, nil)
		}
	}

	self.attrs[index] = val
}

func (self *AttrMap) GetAttr(key *fairy.AttrKey) interface{} {
	index := key.GetIndex()
	if index < len(self.attrs) {
		return self.attrs[index]
	}

	return nil
}

func (self *AttrMap) GetAttrEx(key *fairy.AttrKey, defVal interface{}) interface{} {
	index := key.GetIndex()
	if index < len(self.attrs) {
		return self.attrs[index]
	}

	return defVal
}
