package fairy

/**
 * Identity 用于消息头的序列化，并创建相应Packet,注:无需创建Message
 */
type Identity interface {
	Encode(buffer *Buffer, packet interface{}) error
	Decode(buffer *Buffer) (Packet, error)
}