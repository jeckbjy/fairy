package base

import (
	"sync/atomic"
	"time"

	"github.com/jeckbjy/fairy"
)

const (
	FutureStateNone = iota
	FutureStateSucceed
	FutureStateFailure
	FutureStateCanceled
)

func NewPromise(conn fairy.Conn) *Promise {
	p := &Promise{}
	return p
}

type Promise struct {
	conn  fairy.Conn
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

func (p *Promise) Conn() fairy.Conn {
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
