package fairy

type Future interface {
	IsSuccess() bool
	IsFailure() bool
	IsCanceled() bool
	Conn() Conn
	Cancel() bool
	Wait(msec int64) bool
}

type Promise interface {
	Future
	SetSuccess()
	SetFailure()
}
