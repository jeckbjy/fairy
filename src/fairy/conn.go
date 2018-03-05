package fairy

import "net"

// ConnState
const (
	ConnStateClosed     = 0
	ConnStateOpen       = 1
	ConnStateConnecting = 2
)

// Connection has properties Type ConnId,Uid,OpenId,Data,State,Side and so on
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
	GetState() int32
	SetState(state int32)
	IsServerSide() bool
	IsClientSide() bool
	GetTransport() Transport
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	// operations
	Close()
	Flush()
	Read() *Buffer
	Write(buffer *Buffer)
	Send(obj interface{})
}