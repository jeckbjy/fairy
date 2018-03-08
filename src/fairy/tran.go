package fairy

var (
	// CfgReconnectOpen 设置是否开启自动重连
	CfgReconnectOpen = NewAttrKey(AttrCatConfigSystem, "ReconnectOpen")
	// CfgReconnectInterval 设置重连间隔,秒为单位,默认1s
	CfgReconnectInterval = NewAttrKey(AttrCatConfigSystem, "ReconnectInterval")
	// CfgReaderBufferSize 设置读缓冲, 默认1024
	CfgReaderBufferSize = NewAttrKey(AttrCatConfigSystem, "ReaderBufferSize")
)

type Transport interface {
	SetConfig(key *AttrKey, val interface{})
	GetConfig(key *AttrKey) interface{}
	GetFilterChain() FilterChain
	SetFilterChain(chain FilterChain)
	AddFilters(filters ...Filter)
	Listen(host string, kind int) error
	Connect(host string, kind int) (Future, error)
	Reconnect(conn Conn) (Future, error) // 断线重连使用
	Start()
	Stop()
	Wait()
}
