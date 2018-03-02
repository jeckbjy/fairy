package base

import (
	"fairy"
	"sync/atomic"
	"time"
)

func NewFuture() *BaseFuture {
	future := &BaseFuture{}
	future.New()
	return future
}

type BaseFuture struct {
	result int32
}

func (self *BaseFuture) New() {
	self.result = fairy.FUTURE_RESULT_NONE
}

func (self *BaseFuture) Reset() {
	self.result = fairy.FUTURE_RESULT_NONE
}

func (self *BaseFuture) Succeed() bool {
	return atomic.LoadInt32(&self.result) == fairy.FUTURE_RESULT_SUCCEED
}

func (self *BaseFuture) Result() int {
	return int(atomic.LoadInt32(&self.result))
}

func (self *BaseFuture) HasResult() bool {
	return self.IsResult(fairy.FUTURE_RESULT_NONE)
}

func (self *BaseFuture) IsResult(result int) bool {
	return self.Result() == result
}

func (self *BaseFuture) Wait(msec int64) bool {
	self.result = 0
	if msec != -1 {
		start := time.Now().UnixNano() * int64(time.Millisecond)
		for self.IsResult(fairy.FUTURE_RESULT_NONE) {
			time.Sleep(time.Millisecond)
			now := time.Now().UnixNano() * int64(time.Millisecond)
			if now-start >= msec {
				break
			}
		}
	} else {
		for self.IsResult(fairy.FUTURE_RESULT_NONE) {
			time.Sleep(time.Millisecond)
		}
	}

	self.Done(fairy.FUTURE_RESULT_TIMEOUT)
	return self.Succeed()
}

func (self *BaseFuture) Done(result int) {
	atomic.CompareAndSwapInt32(&self.result, fairy.FUTURE_RESULT_NONE, int32(result))
}

func (self *BaseFuture) DoneSucceed() {
	self.Done(fairy.FUTURE_RESULT_SUCCEED)
}

func (self *BaseFuture) DoneFail() {
	self.Done(fairy.FUTURE_RESULT_FAIL)
}

func (self *BaseFuture) DoneTimeout() {
	self.Done(fairy.FUTURE_RESULT_TIMEOUT)
}

/////////////////////////////////////////////////////////////////////////////////////
// BaseConnectFuture
/////////////////////////////////////////////////////////////////////////////////////

func NewConnectFuture(conn fairy.Connection) *BaseConnectFuture {
	future := &BaseConnectFuture{}
	future.BaseFuture.New()
	future.conn = conn
	return future
}

type BaseConnectFuture struct {
	BaseFuture
	conn fairy.Connection
}

func (self *BaseConnectFuture) Get(msec int64) (fairy.Connection, bool) {
	self.Wait(msec)
	return self.conn, self.Succeed()
}
