package fairy

import "net"

// IConn asynchronous connection,
type IConn interface {
	AttrMap
	GetId() uint
	GetTag() string
	SetTag(tag string)
	GetData() interface{}
	SetData(data interface{})
	IsActive() bool             // 连接是否正常
	IsConnector() bool          // 是否通过调用Connect产生,否则Listen产生
	LocalAddr() net.Addr        // 本地地址
	RemoteAddr() net.Addr       // 远程地址
	Close()                     // 异步关闭,会等待数据发送完
	Read() *Buffer              // 异步读缓存
	Write(buffer *Buffer) error // 异步写数据
	Send(msg interface{}) error // 异步发消息
}

// ITran Transport,用于创建Connection
type ITran interface {
	GetChain() IFilterChain
	SetOptions(option ...Option)
	AddFilters(filters ...IFilter)
	Listen(host string, options ...Option) error  // 支持TagOption
	Connect(host string, options ...Option) error // 支持TagOption,SyncOption,ReconnectOption
	Start()
	Stop()
}
