package soa

import (
	"github.com/jeckbjy/fairy"
)

// Info 节点信息
type Info struct {
	ID          uint64                 // 唯一标识
	InnerAddr   string                 // 内网IP地址
	OuterAddr   string                 // 外网IP地址
	Port        uint                   // 服务端口,是否需要支持多个?
	Load        uint                   // 负载值
	PubSerivces []string               // 提供的服务
	Data        map[string]interface{} // TODO:自定义数据
	conn        fairy.Conn             // 连接Conn
}

// GetConn 返回关联的Conn
func (info *Info) GetConn() fairy.Conn {
	return info.conn
}

// SetConn 设置Conn
func (info *Info) SetConn(conn fairy.Conn) {
	info.conn = conn
}

// SetData 设置关联数据
func (info *Info) SetData(key string, val interface{}) {
}

// InfoEx 提供自己订阅的信息,但下发时并不需要告诉他人
type InfoEx struct {
	Info
	SubServices []string // 订阅的服务
}
