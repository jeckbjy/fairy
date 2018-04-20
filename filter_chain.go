package fairy

// FilterChain 责任链接口
type FilterChain interface {
	Len() int
	AddFirst(filter Filter)
	AddLast(filter Filter)
	HandleOpen(conn Conn)
	HandleClose(conn Conn)
	HandleRead(conn Conn)
	HandleWrite(conn Conn, msg interface{})
	HandleError(conn Conn, err error)
}
