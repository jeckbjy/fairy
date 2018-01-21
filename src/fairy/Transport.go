package fairy

// 断线重连处理??
type Transport interface {
	SetFilterChain(chain FilterChain)
	AddFilters(filters ...Filter)
	Listen(host string, ctype int)
	Connect(host string, ctype int) ConnectFuture
	Start()
	Stop()
	Wait()
}
