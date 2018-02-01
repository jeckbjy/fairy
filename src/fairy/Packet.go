package fairy

// 消息调用结果
const (
	PacketResultNone    = 0
	PacketResultOk      = 1
	PacketResultFail    = 2
	PacketResultTimeout = 3
)

/**
 * packet 包含消息头和消息体两部分,配合Identity和Codec使用
 * 消息头: 接口定义了一些常用的消息头,方便直接使用，无需转型，
 * 		  id,name:id和name用于标识Packet类型，id为0时由name标识
 *		  result,serialid,time:通常用于rpc系统,serialid用于找到唯一回调函数,result标识结果,time用于计算超时
 * 		  checksum:用于校验消息包的完整性
 * 消息体: 由Message表示，可以是string,json,protobu编码，具体编解码由Codec实现
 */
type Packet interface {
	GetId() uint
	SetId(id uint)
	GetName() string
	SetName(name string)
	GetResult() uint
	SetResult(result uint)
	GetSerialId() uint
	SetSerialId(id uint)
	GetTime() uint
	SetTime(time uint)
	GetChecksum() uint64
	SetChecksum(val uint64)
	GetMessage() interface{}
	SetMessage(msg interface{})
}
