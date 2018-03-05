package fairy

var (
	KeyReconnectInterval = NewAttrKey(AttrCatConfigSystem, "ReconnectInterval")
	KeyReaderBufferSize  = NewAttrKey(AttrCatConfigSystem, "ReaderBufferSize")
)

type Transport interface {
	SetConfig(key *AttrKey, val interface{})
	GetConfig(key *AttrKey) interface{}
	SetFilterChain(chain FilterChain)
	AddFilters(filters ...Filter)
	Listen(host string, kind int) error
	Connect(host string, kind int) (ConnectFuture, error)
	Start()
	Stop()
	Wait()
}
