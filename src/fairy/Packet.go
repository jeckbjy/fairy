package fairy

const (
	PACKET_RESULT_NONE = iota
	PACKET_RESULT_OK
	PACKET_RESULT_FAIL
	PACKET_RESULT_TIMEOUT
)

/*
packet=head+body
id,name:用于标识消息体类型
message:真正的消息体
result,serialid,time:可用于标识rpc结果，唯一ID和超时时间
*/
type Packet interface {
	GetResult() uint
	SetResult(result uint)
	GetSerialId() uint
	SetSerialId(id uint)
	GetTime() uint
	SetTime(time uint)
	GetId() uint
	SetId(id uint)
	GetName() string
	SetName(name string)
	GetMessage() interface{}
	SetMessage(msg interface{})
}
