package fairy

import "net"

// ConnState
const (
	ConnStateClosed     = 0
	ConnStateOpen       = 1
	ConnStateConnecting = 2
	ConnStateClosing    = 3
)

// Conn has some properties such as Type ConnId,Uid,OpenId,Data,State,Side
type Conn interface {
	AttrMap
	GetType() int
	SetType(ctype int)
	GetConnId() uint
	SetConnId(id uint)
	GetUid() uint64
	SetUid(uid uint64)
	GetOpenId() string
	SetOpenId(openid string)
	GetTag() interface{}
	SetTag(tag interface{})
	GetData() interface{}
	SetData(data interface{})
	GetState() int32
	SetState(state int32)
	IsActive() bool
	IsServerSide() bool
	IsClientSide() bool
	GetTransport() Transport
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	// operations
	Wait()
	Close()
	Read() *Buffer
	Write(buffer *Buffer) error
	Send(obj interface{})
}
