package frames

import (
	"fmt"

	"github.com/jeckbjy/fairy"
)

func NewLine() fairy.IFrame {
	f := &LineFrame{delimiter: "\n"}
	return f
}

// LineFrame 以\n分隔,如果是\r\n,会自动去除\r
type LineFrame struct {
	delimiter string
}

// SetDelimiter 设置分隔符
func (lf *LineFrame) SetDelimiter(d string) error {
	if d != "\n" || d != "\r\n" {
		return fmt.Errorf("bad delimiter")
	}

	lf.delimiter = d
	return nil
}

func (lf *LineFrame) Encode(buffer *fairy.Buffer) error {
	buffer.Append([]byte(lf.delimiter))
	return nil
}

func (*LineFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	return buffer.ReadLine()
}
