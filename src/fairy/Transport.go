package fairy

var (
	KeyReconnectInterval = NewAttrKey(AttrCatConfigSystem, "ReconnectInterval")
	KeyReaderBufferSize  = NewAttrKey(AttrCatConfigSystem, "ReaderBufferSize")
)

type Transport interface {
	SetConfig(key *AttrKey, val string)
	GetConfig(key *AttrKey) interface{}
	SetFilterChain(chain FilterChain)
	AddFilters(filters ...Filter)
	Listen(host string, ctype int)
	Connect(host string, ctype int) ConnectFuture
	Start()
	Stop()
	Wait()
}
