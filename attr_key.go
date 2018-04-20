package fairy

import (
	"fmt"
	"sync"
)

const (
	// AttrKindConf 配置使用
	AttrKindConf = "conf"
	// AttrKindConn Conn中存储数据
	AttrKindConn = "conn"
	// AttrKindCtx filterContext中使用
	AttrKindCtx = "context"
)

// NewAttrKey 创建NewAttrKey
func NewAttrKey(kind string, name string) *AttrKey {
	return gAttrMgr.Create(kind, name)
}

// AttrKey 将string延迟映射到唯一索引,不同进程间,索引并不一定一样
type AttrKey struct {
	kind  string
	name  string
	index int
}

// Index return attr index
func (attr *AttrKey) Index() int {
	if attr.index == -1 {
		gAttrMgr.Generate(attr)
	}

	return attr.index
}

// Name return attr name
func (attr *AttrKey) Name() string {
	return attr.name
}

// String return attr stringify
func (attr *AttrKey) String() string {
	return fmt.Sprintf(":%s:%s:%v", attr.kind, attr.name, attr.index)
}

//////////////////////////////////////////////////////
// AttrMgr:global AttrKey manager
//////////////////////////////////////////////////////
var gAttrMgr = &zAttrMgr{attrs: make(map[string]*AttrKey), indexs: make(map[string]int)}

type zAttrMgr struct {
	attrs  map[string]*AttrKey
	indexs map[string]int
	mutex  sync.Mutex
}

// Create 查找或创建一个AttrKey
func (mgr *zAttrMgr) Create(kind string, name string) *AttrKey {
	var attr *AttrKey
	mgr.mutex.Lock()
	key := fmt.Sprintf("%s:%s", kind, name)
	if val, ok := mgr.attrs[key]; ok {
		attr = val
	} else {
		attr = &AttrKey{kind: kind, name: name, index: -1}
		mgr.attrs[key] = attr
	}

	mgr.mutex.Unlock()
	return attr
}

// Generate 生成一个索引,从零开始
func (mgr *zAttrMgr) Generate(attr *AttrKey) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	if attr.index != -1 {
		return
	}

	kind := attr.kind
	index := 0
	if val, ok := mgr.indexs[kind]; ok {
		index = val
	}

	attr.index = index

	mgr.indexs[kind] = index + 1
}
