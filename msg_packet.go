package fairy

const (
	PacketResultUnknown = 0
	PacketResultSuccess = 1
	PacketResultFailure = 2
	PacketResultTimeout = 3
)

// Packet 用于定义消息包，分为消息头和消息体
// 消息头:接口定义了一些常用的消息头,方便直接使用，无需向下转型
//   	1:id,name:id和name用于标识Packet类型，id为0时由name标识,通常必须存在
//		2:rpcid,result,用于库默认提供的rpc系统,根据配置，不一定有用
// 		3:checksum并没有放到packet中定义，因为通常有frame来校验完整性
// 消息体:由Message表示，可以是string,json,protobu编码，具体编解码由Codec实现
type Packet interface {
	Reset()
	GetId() uint
	SetId(id uint)
	GetName() string
	SetName(name string)
	GetMessage() interface{}
	SetMessage(msg interface{})
	GetRpcId() uint64
	SetRpcId(id uint64)
	SetResult(result uint)
	GetResult() uint
	// 以下仅仅是result的封装
	SetTimeout()
	SetSuccess()
	SetFailure()
	IsTimeout() bool
	IsSuccess() bool
	IsFailure() bool
}

// Identity 用于消息头的序列化，并创建相应Packet,注:无需创建Message
type Identity interface {
	Encode(buffer *Buffer, packet interface{}) error
	Decode(buffer *Buffer) (Packet, error)
}

// Codec 用于消息的序列化
type Codec interface {
	Encode(msg interface{}, buffer *Buffer) error
	Decode(msg interface{}, buffer *Buffer) error
}

// Frame 帧序列
type Frame interface {
	Encode(buffer *Buffer) error
	Decode(buffer *Buffer) (*Buffer, error)
}
