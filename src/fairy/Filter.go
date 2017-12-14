package fairy

// 每个Filter要求是无状态的，如果需要存储数据，请使用AttrMap
type Filter interface {
	HandleRead(ctx FilterContext) FilterAction
	HandleWrite(ctx FilterContext) FilterAction
	HandleOpen(ctx FilterContext) FilterAction
	HandleClose(ctx FilterContext) FilterAction
	HandleError(ctx FilterContext) FilterAction
	//NotifyUpstram(ctx FilterContext) FilterAction
	//NotifyDownstram(ctx FilterContext) FilterAction
}
