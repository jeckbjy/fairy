package snet

import (
	"fairy/base"
	"net"
	"sync"
)

type StreamTran struct {
	base.Tran
	listeners []net.Listener
	wg        sync.WaitGroup
	stopped   bool
}

func (t *StreamTran) Create() {
	t.SetDefaultConfig()
	t.stopped = false
}

func (t *StreamTran) IsStopped() bool {
	return t.stopped
}

func (t *StreamTran) AddListener(l net.Listener) {
	t.listeners = append(t.listeners, l)
}

func (t *StreamTran) AddGroup() {
	t.wg.Add(1)
}

func (t *StreamTran) Done() {
	t.wg.Done()
}

func (t *StreamTran) Wait() {
	t.wg.Wait()
}

func (t *StreamTran) Stop() {
	t.stopped = true
	// close all listener
	for _, listener := range t.listeners {
		listener.Close()
	}

	t.listeners = nil
}

func (t *StreamTran) OnExit() {
	t.Stop()
}
