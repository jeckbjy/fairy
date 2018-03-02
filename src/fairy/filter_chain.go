package fairy

type FilterChain interface {
	Len() int
	AddFirst(filter Filter)
	AddLast(filter Filter)
	HandleOpen(conn Connection)
	HandleClose(conn Connection)
	HandleRead(conn Connection)
	HandleWrite(conn Connection, msg interface{})
	HandleError(conn Connection, err error)
}
