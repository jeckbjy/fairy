package fairy

const (
	FUTURE_RESULT_NONE    = 0
	FUTURE_RESULT_SUCCEED = 1
	FUTURE_RESULT_FAIL    = 2
	FUTURE_RESULT_TIMEOUT = 3
)

// todo:
// CloseFuture,WriteFuture
type Future interface {
	Succeed() bool
	Result() int
	Wait(msec int64) bool
	Done(result int)
}

type ConnectFuture interface {
	Future
	Get(msec int64) (Connection, bool)
}
