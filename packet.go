package fairy

// IPacket 用于定义消息包，分为消息头和消息体
// 消息头:Id,Name,用于唯一标识消息,Id为0时使用Name
// 消息体:可以是string,json,protobu编码，具体编解码由Codec实现
type IPacket interface {
	GetId() uint
	SetId(id uint)
	GetName() string
	SetName(name string)
	GetMessage() interface{}
	SetMessage(msg interface{})
}

// IIdentity 用于消息头的序列化，并创建相应Packet,注:无需创建Message
type IIdentity interface {
	Encode(buffer *Buffer, pkt IPacket) error
	Decode(buffer *Buffer) (IPacket, error)
}

// ICodec 用于消息的序列化
type ICodec interface {
	Encode(buffer *Buffer, msg interface{}) error
	Decode(buffer *Buffer, msg interface{}) error
}

// IFrame 帧序列,用于粘包处理
type IFrame interface {
	Encode(buffer *Buffer) error
	Decode(buffer *Buffer) (*Buffer, error)
}
