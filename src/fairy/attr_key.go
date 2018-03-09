package fairy

import "sync"

const (
	// 配置使用
	AttrKindConfig = iota
	// conn长期持有使用
	AttrKindConn
	// filter之间传递数据使用
	AttrKindFilter
)

// NewAttrKey lazy create attr key
func NewAttrKey(kind int, name string) *AttrKey {
	return findOrCreate(kind, name)
}

// AttrKey string to int index, lazy init index
type AttrKey struct {
	name  string
	index int
	owner *attrKeyMap
}

// GetIndex lazy create index and return index
func (ak *AttrKey) GetIndex() int {
	if ak.index == -1 {
		ak.owner.newIndex(ak)
	}

	return ak.index
}

// GetName return attr name
func (ak *AttrKey) GetName() string {
	return ak.name
}

//////////////////////////////////////////////////////
// AttrKeyMap:global manager
//////////////////////////////////////////////////////
var gAttrKeyMapArray []*attrKeyMap
var gAttrKeyMux sync.Mutex

func findOrCreate(category int, name string) *AttrKey {
	gAttrKeyMux.Lock()
	defer gAttrKeyMux.Unlock()

	if category >= len(gAttrKeyMapArray) {
		count := category - len(gAttrKeyMapArray) + 1
		gAttrKeyMapArray = append(gAttrKeyMapArray, make([]*attrKeyMap, count)...)
	}

	if gAttrKeyMapArray[category] == nil {
		gAttrKeyMapArray[category] = newAttrKeyMap(category)
	}

	return gAttrKeyMapArray[category].getOrCreate(name)
}

func newAttrKeyMap(category int) *attrKeyMap {
	akm := &attrKeyMap{category: category, index: 0}
	akm.attrs = make(map[string]*AttrKey)
	return akm
}

type attrKeyMap struct {
	category int
	index    int
	attrs    map[string]*AttrKey
}

func (akm *attrKeyMap) getOrCreate(name string) *AttrKey {
	attr, ok := akm.attrs[name]
	if !ok {
		attr = &AttrKey{name: name, index: -1, owner: akm}
		akm.attrs[name] = attr
	}

	return attr
}

func (akm *attrKeyMap) newIndex(attr *AttrKey) {
	gAttrKeyMux.Lock()
	attr.index = akm.index
	akm.index++
	gAttrKeyMux.Unlock()
}
