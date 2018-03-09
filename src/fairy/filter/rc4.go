package filter

import (
	"crypto/rc4"
	"fairy"
	"fairy/base"
	"fmt"
)

var (
	// KeyRC4 attr key in conn
	KeyRC4 = fairy.NewAttrKey(fairy.AttrKindConn, "rc4")
)

// NewRC4 create rc4 filter with secret default key
func NewRC4(key string, useConnKey bool) fairy.Filter {
	f := &rc4Filter{useConnKey: useConnKey}
	if key != "" {
		f.key = []byte(key)
	}

	return f
}

// if no set key, won't crypt
type rc4Filter struct {
	base.BaseFilter
	key        []byte // default key
	useConnKey bool
}

func (rf *rc4Filter) newCipher(conn fairy.Conn) (*rc4.Cipher, error) {
	key := rf.key

	if rf.useConnKey {
		attr := conn.GetAttr(KeyRC4)
		if attr != nil {
			switch attr.(type) {
			case string:
				key = []byte(attr.(string))
			case []byte:
				key = attr.([]byte)
			default:
				return nil, fmt.Errorf("rc4 key must be string or []byte")
			}
		}
	}

	// if no key, do not crypto
	if key == nil {
		return nil, nil
	}

	return rc4.NewCipher(key)
}

func (rf *rc4Filter) handle(ctx fairy.FilterContext) fairy.FilterAction {
	buf, ok := ctx.GetMessage().(*fairy.Buffer)
	if !ok {
		return ctx.GetNextAction()
	}

	cipher, err := rf.newCipher(ctx.GetConn())
	if err != nil {
		return ctx.ThrowError(err)
	}

	if cipher != nil {
		buf.Visit(func(data []byte) bool {
			cipher.XORKeyStream(data, data)
			return true
		})
	}

	return ctx.GetNextAction()
}

func (rf *rc4Filter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	return rf.handle(ctx)
}

func (rf *rc4Filter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	return rf.handle(ctx)
}
