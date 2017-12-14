package fairy

// 消息ID标识:Length+[Identity(head)+Message(body)]
// IdentityCodec
type Identity interface {
	Decode(buffer *Buffer) (Packet, error)
	Encode(buffer *Buffer, packet interface{}) error
}
