package fairy

// IFilter must be stateless,if need data, can get from Conn or FilteCtx
type IFilter interface {
	Name() string               // unique name for filter
	HandleRead(ctx IFilterCtx)  // read data
	HandleWrite(ctx IFilterCtx) // send data
	HandleOpen(ctx IFilterCtx)  // open by connect or listen
	HandleClose(ctx IFilterCtx)
	HandleError(ctx IFilterCtx)
}

// IFilterCtx filter上下文
type IFilterCtx interface {
	AttrMap
	GetConn() IConn
	SetData(data interface{}) // 用于Filter间透传数据
	GetData() interface{}     // 获取数据
	Error(err error)          // 抛出错误
	Next()                    // 执行下一个
	Jump(index int) error     // 跳转到某个索引,可以负索引
	JumpBy(name string) error // 通过名字跳转
}

// IFilterChain 递归执行每一个Filter
type IFilterChain interface {
	Len() int
	IndexOf(name string) int     // 通过名字查询索引
	AddFirst(filters ...IFilter) // 前边插入,Prepend
	AddLast(filters ...IFilter)  // 后边插入,Append
	HandleOpen(conn IConn)
	HandleClose(conn IConn)
	HandleRead(conn IConn)
	HandleWrite(conn IConn, msg interface{})
	HandleError(conn IConn, err error)
}
