package fairy

var (
	// CfgReconnectOpen 设置是否开启自动重连
	CfgReconnectOpen = NewAttrKey(AttrKindConf, "ReconnectOpen")
	// CfgReconnectInterval 设置重连间隔,秒为单位,默认1s
	CfgReconnectInterval = NewAttrKey(AttrKindConf, "ReconnectInterval")
	// CfgReaderBufferSize 设置读缓冲, 默认1024
	CfgReaderBufferSize = NewAttrKey(AttrKindConf, "ReaderBufferSize")
	// CfgAutoRead 自动开启读协程,false,socket将不会读数据
	CfgAutoRead = NewAttrKey(AttrKindConf, "AutoRead")
)

// Tran 负责Conn的创建
type Tran interface {
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
