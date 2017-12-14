package fairy

// Packet=head+body
// 基础消息包:包含名字，唯一ID标识，和消息体
type Packet interface {
	GetId() uint
	GetName() string
	GetMessage() interface{}
	SetId(id uint)
	SetName(name string)
	SetMessage(msg interface{})
}
