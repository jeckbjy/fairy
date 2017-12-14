package fairy

import "errors"

var (
	ErrConnectFail = errors.New("err connect fail!")
	ErrReadFail    = errors.New("err read fail!")
)
