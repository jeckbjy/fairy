package fairy

type AttrKey struct {
	name  string
	index int
}

func (self *AttrKey) GetIndex() int {
	return self.index
}

func (self *AttrKey) GetName() string {
	return self.name
}

func NewAttrKey(category int, name string) *AttrKey {
	attr_map := GetAttrKeyMap(category)
	return attr_map.Create(name)
}

//////////////////////////////////////////////////////
// AttrKeyMap:global manager
//////////////////////////////////////////////////////
var gAttrKeyMapArray []*AttrKeyMap

func GetAttrKeyMap(category int) *AttrKeyMap {
	count := category - len(gAttrKeyMapArray) + 1
	if count > 0 {
		for i := 0; i < count; i++ {
			gAttrKeyMapArray = append(gAttrKeyMapArray)
		}
	}

	if gAttrKeyMapArray[category] == nil {
		gAttrKeyMapArray[category] = NewAttrKeyMap(category)
	}

	return gAttrKeyMapArray[category]
}

func NewAttrKeyMap(category int) *AttrKeyMap {
	attr_key_map := &AttrKeyMap{category: category, index: 0}
	attr_key_map.attrs = make(map[string]*AttrKey)
	return attr_key_map
}

type AttrKeyMap struct {
	category int
	index    int
	attrs    map[string]*AttrKey
}

func (self *AttrKeyMap) Create(name string) *AttrKey {
	attr, ok := self.attrs[name]
	if !ok {
		attr = &AttrKey{name: name, index: self.index}
		self.attrs[name] = attr
		self.index++
	}

	return attr
}
