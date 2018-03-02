package fairy

type AttrMap interface {
	HasAttr(key *AttrKey) bool
	SetAttr(key *AttrKey, val interface{})
	GetAttr(key *AttrKey) interface{}
	GetAttrEx(key *AttrKey, defVal interface{}) interface{}
}
