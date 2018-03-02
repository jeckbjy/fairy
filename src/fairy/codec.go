package fairy

// MessageCodec
type Codec interface {
	Encode(msg interface{}, buffer *Buffer) error
	Decode(msg interface{}, buffer *Buffer) error
}
