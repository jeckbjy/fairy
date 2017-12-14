package fairy

type FilterChain interface {
	AddFirst(filter Filter)
	AddLast(filter Filter)
	HandleOpen(conn Connection)
	HandleClose(conn Connection)
	HandleRead(conn Connection)
	HandleWrite(conn Connection, msg interface{})
	HandleError(conn Connection, err error)
}
