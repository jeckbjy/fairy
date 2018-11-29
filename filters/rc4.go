package filters

import (
	"crypto/rc4"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

// NewRC4 create rc4 filter with secret default key
func NewRC4(key string) *RC4Filter {
	return &RC4Filter{key: []byte(key)}
}

// RC4Filter rc4加密算法,key固定不可变
type RC4Filter struct {
	base.Filter
	key []byte // default key
}

func (rf *RC4Filter) Name() string {
	return "RC4Filter"
}

func (rf *RC4Filter) SetKey(key string) {
	rf.key = []byte(key)
}

func (rf *RC4Filter) HandleRead(ctx fairy.IFilterCtx) {
	rf.handle(ctx)
}

func (rf *RC4Filter) HandleWrite(ctx fairy.IFilterCtx) {
	rf.handle(ctx)
}

func (rf *RC4Filter) handle(ctx fairy.IFilterCtx) {
	buf, ok := ctx.GetData().(*fairy.Buffer)
	if !ok {
		ctx.Next()
		return
	}

	cipher, err := rc4.NewCipher(rf.key)
	if err != nil {
		// warning
		return
	}

	if cipher != nil {
		buf.Visit(func(data []byte) bool {
			cipher.XORKeyStream(data, data)
			return true
		})
	}

	ctx.Next()
}
