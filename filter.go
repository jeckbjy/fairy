package fairy

// FilterAction 用于
type FilterAction interface {
	Type() int
}

// Filter Stateless标识是否有状态
type Filter interface {
	Stateless() bool
	HandleRead(ctx FilterContext) FilterAction
	HandleWrite(ctx FilterContext) FilterAction
	HandleOpen(ctx FilterContext) FilterAction
	HandleClose(ctx FilterContext) FilterAction
	HandleError(ctx FilterContext) FilterAction
	//NotifyUpstram(ctx FilterContext) FilterAction
	//NotifyDownstram(ctx FilterContext) FilterAction
}
