package soa

// 注册提供的服务和订阅服务
type soaRegisterReq struct {
	Info *InfoEx
}

// 分发订阅的服务
type soaRegisterRsp struct {
	SubInfos []*Info
}

// 通知节点删除
type soaRemoveMsg struct {
	ID uint64
}

// 更新Load
type soaUpdateLoadMsg struct {
	Load uint
}
