package fairy

var (
	CfgReconnectInterval = NewAttrKey(AttrCatConfigSystem, "ReconnectInterval")
	CfgReconnectCount    = NewAttrKey(AttrCatConfigSystem, "ReconnectCount")
	CfgReaderBufferSize  = NewAttrKey(AttrCatConfigSystem, "ReaderBufferSize")
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
