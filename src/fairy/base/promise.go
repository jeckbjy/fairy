package base

import (
	"fairy"
	"sync/atomic"
	"time"
)

const (
	FutureStateNone = iota
	FutureStateSucceed
	FutureStateFailure
	FutureStateCanceled
)

func NewPromise(conn fairy.Connection) *Promise {
	p := &Promise{}
	return p
}

type Promise struct {
	conn  fairy.Connection
	state int32
}

func (p *Promise) isState(s int32) bool {
	return atomic.LoadInt32(&p.state) == s
}

func (p *Promise) setState(s int32) bool {
	return atomic.CompareAndSwapInt32(&p.state, FutureStateNone, s)
}

func (p *Promise) IsSuccess() bool {
	return p.isState(FutureStateSucceed)
}

func (p *Promise) IsFailure() bool {
	return p.isState(FutureStateFailure)
}

func (p *Promise) IsCanceled() bool {
	return p.isState(FutureStateCanceled)
}

func (p *Promise) Conn() fairy.Connection {
	return p.conn
}

func (p *Promise) Cancel() bool {
	return p.setState(FutureStateCanceled)
}

func (p *Promise) Wait(msec int64) bool {
	p.state = FutureStateNone
	if msec != -1 {
		start := time.Now().UnixNano() * int64(time.Millisecond)
		for p.isState(FutureStateNone) {
			time.Sleep(time.Millisecond)
			now := time.Now().UnixNano() * int64(time.Millisecond)
			if now-start >= msec {
				break
			}
		}
	} else {
		for p.isState(FutureStateNone) {
			time.Sleep(time.Millisecond)
		}
	}

	p.setState(FutureStateCanceled)
	return p.IsSuccess()
}

func (p *Promise) SetSuccess() {
	p.setState(FutureStateSucceed)
}

func (p *Promise) SetFailure() {
	p.setState(FutureStateFailure)
}

// func NewFuture() *BaseFuture {
// 	future := &BaseFuture{}
// 	future.New()
// 	return future
// }

// type BaseFuture struct {
// 	result int32
// }

// func (self *BaseFuture) New() {
// 	self.result = fairy.FUTURE_RESULT_NONE
// }

// func (self *BaseFuture) Reset() {
// 	self.result = fairy.FUTURE_RESULT_NONE
// }

// func (self *BaseFuture) Succeed() bool {
// 	return atomic.LoadInt32(&self.result) == fairy.FUTURE_RESULT_SUCCEED
// }

// func (self *BaseFuture) Result() int {
// 	return int(atomic.LoadInt32(&self.result))
// }

// func (self *BaseFuture) HasResult() bool {
// 	return self.IsResult(fairy.FUTURE_RESULT_NONE)
// }

// func (self *BaseFuture) IsResult(result int) bool {
// 	return self.Result() == result
// }

// func (self *BaseFuture) Wait(msec int64) bool {
// 	self.result = 0
// 	if msec != -1 {
// 		start := time.Now().UnixNano() * int64(time.Millisecond)
// 		for self.IsResult(fairy.FUTURE_RESULT_NONE) {
// 			time.Sleep(time.Millisecond)
// 			now := time.Now().UnixNano() * int64(time.Millisecond)
// 			if now-start >= msec {
// 				break
// 			}
// 		}
// 	} else {
// 		for self.IsResult(fairy.FUTURE_RESULT_NONE) {
// 			time.Sleep(time.Millisecond)
// 		}
// 	}

// 	self.Done(fairy.FUTURE_RESULT_TIMEOUT)
// 	return self.Succeed()
// }

// func (self *BaseFuture) Done(result int) {
// 	atomic.CompareAndSwapInt32(&self.result, fairy.FUTURE_RESULT_NONE, int32(result))
// }

// func (self *BaseFuture) DoneSucceed() {
// 	self.Done(fairy.FUTURE_RESULT_SUCCEED)
// }

// func (self *BaseFuture) DoneFail() {
// 	self.Done(fairy.FUTURE_RESULT_FAIL)
// }

// func (self *BaseFuture) DoneTimeout() {
// 	self.Done(fairy.FUTURE_RESULT_TIMEOUT)
// }

// /////////////////////////////////////////////////////////////////////////////////////
// // BaseConnectFuture
// /////////////////////////////////////////////////////////////////////////////////////

// func NewConnectFuture(conn fairy.Connection) *BaseConnectFuture {
// 	future := &BaseConnectFuture{}
// 	future.BaseFuture.New()
// 	future.conn = conn
// 	return future
// }

// type BaseConnectFuture struct {
// 	BaseFuture
// 	conn fairy.Connection
// }

// func (self *BaseConnectFuture) Get(msec int64) (fairy.Connection, bool) {
// 	self.Wait(msec)
// 	return self.conn, self.Succeed()
// }
