package fairy

var (
	// CfgReconnectOpen 设置是否开启自动重连
	CfgReconnectOpen = NewAttrKey(AttrKindConfig, "ReconnectOpen")
	// CfgReconnectInterval 设置重连间隔,秒为单位,默认1s
	CfgReconnectInterval = NewAttrKey(AttrKindConfig, "ReconnectInterval")
	// CfgReaderBufferSize 设置读缓冲, 默认1024
	CfgReaderBufferSize = NewAttrKey(AttrKindConfig, "ReaderBufferSize")
	// CfgAutoRead 自动开启读协程,false,socket将不会读数据
	CfgAutoRead = NewAttrKey(AttrKindConfig, "AutoRead")
)

// Transport 负责Conn的创建
type Transport interface {
	SetConfig(key *AttrKey, val interface{})
	GetConfig(key *AttrKey) interface{}
	GetFilterChain() FilterChain
	SetFilterChain(chain FilterChain)
	AddFilters(filters ...Filter)
	Listen(host string, tag interface{}) error
	Connect(host string, tag interface{}) (Future, error)
	Reconnect(conn Conn) (Future, error) // 断线重连使用
	Start()
	Stop()
	Wait()
}
