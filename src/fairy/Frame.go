package fairy

// 帧序列，FrameCodec
type Frame interface {
	Encode(buffer *Buffer) error
	Decode(buffer *Buffer) (*Buffer, error)
}
