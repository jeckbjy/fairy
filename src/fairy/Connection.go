package fairy

import "net"

const (
	CONN_STATE_CLOSED     = 0
	CONN_STATE_OPEN       = 1
	CONN_STATE_CONNECTING = 2
	CONN_STATE_CLOSING    = 3
)

// Future support(CloseFuture,WriteFuture)???
// 常用属性：Type,ConnId,Uid，OpenId，UserData
type Connection interface {
	AttrMap
	GetType() int
	SetType(ctype int)
	GetConnId() uint
	SetConnId(id uint)
	GetUid() uint64
	SetUid(uid uint64)
	GetOpenId() string
	SetOpenId(openid string)
	GetData() interface{}
	SetData(data interface{})
	GetState() int
	IsState(state int) bool
	IsServerSide() bool
	IsClientSide() bool
	GetTransport() Transport
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() Future
	Flush()
	Write(buffer *Buffer)
	Read() *Buffer
}
