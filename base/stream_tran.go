package base

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/jeckbjy/fairy"
)

// NewTran create stream transport
func NewTran(itran IStreamTran) fairy.ITran {
	tran := &StreamTran{}
	tran.init(itran)
	return tran
}

var syncReconnectMax = 10
var errBadTran = errors.New("bad transport,is nil")

type OnAccept func(conn net.Conn, err error)

type IStreamTran interface {
	Connect(host string, options ...fairy.Option) (net.Conn, error)
	Listen(host string, options ...fairy.Option) (net.Listener, error) // 监听端口
	Serve(l net.Listener, cb OnAccept)                                 // 循环Accept
}

// StreamTran 面向长连接Transport
type StreamTran struct {
	Tran
	itran                IStreamTran
	listeners            []net.Listener
	wg                   sync.WaitGroup
	running              bool
	cfgReconnectCount    int // 自动重连次数,-1标识无限制
	cfgReconnectInterval int // 自动重连间隔,单位是秒
}

func (st *StreamTran) init(tran IStreamTran) {
	st.running = true
	st.itran = tran
	st.cfgReconnectCount = -1
	st.cfgReconnectInterval = 1
}

func (st *StreamTran) SetOptions(options ...fairy.Option) {
	for _, op := range options {
		switch op.(type) {
		case *fairy.ReconnectOption:
			rop := op.(*fairy.ReconnectOption)
			st.cfgReconnectCount = rop.Count
			if rop.Count != 0 {
				st.cfgReconnectInterval = rop.Interval
			}
		}
	}
}

// Listen 监听连接
func (st *StreamTran) Listen(host string, options ...fairy.Option) error {
	if st.itran == nil {
		return errBadTran
	}

	listener, err := st.itran.Listen(host)
	if err != nil {
		return err
	}

	st.listeners = append(st.listeners, listener)
	st.wg.Add(1)

	// check tag
	tag := ""
	for _, op := range options {
		if to, ok := op.(*fairy.TagOption); ok {
			tag = to.Tag
		}
	}

	go func() {
		st.itran.Serve(listener, func(conn net.Conn, err error) {
			if err != nil {
				return
			}
			newConn := NewConn(st, false, tag)
			newConn.Open(conn)
		})
		st.wg.Done()
	}()

	return nil
}

// Connect 异步连接
func (st *StreamTran) Connect(host string, options ...fairy.Option) error {
	if st.itran == nil {
		return errBadTran
	}

	tag := ""
	sync := false
	count := 0
	interval := -1
	for _, op := range options {
		switch op.(type) {
		case *fairy.TagOption:
			tag = op.(*fairy.TagOption).Tag
		case *fairy.SyncOption:
			sync = op.(*fairy.SyncOption).Flag
		case *fairy.ReconnectOption:
			ro := op.(*fairy.ReconnectOption)
			count = ro.Count
			interval = ro.Interval
		}
	}

	if interval == -1 {
		interval = st.cfgReconnectInterval
	}

	newConn := NewConn(st, true, tag)
	newConn.SetHost(host)

	if sync {
		// 同步连接,防止连接失败导致永久阻塞
		if count < 0 {
			count = syncReconnectMax
		}

		return st.tryConnect(newConn, host, count, interval)
	} else {
		// 异步连接,默认
		st.wg.Add(1)
		go func() {
			st.tryConnect(newConn, host, count, interval)
			st.wg.Done()
		}()

		return nil
	}
}

// Reconnect 断线重连
func (st *StreamTran) Reconnect(sconn *StreamConn) {
	// 无需断线重连
	if st.cfgReconnectCount == 0 {
		return
	}

	st.wg.Add(1)
	go func() {
		st.tryConnect(sconn, sconn.GetHost(), st.cfgReconnectCount, st.cfgReconnectInterval)
		st.wg.Done()
	}()
}

// 尝试连接
func (st *StreamTran) tryConnect(newConn *StreamConn, host string, cfgCount, cfgInterval int) error {
	count := 0
	for {
		count++
		conn, err := st.itran.Connect(host)
		if err == nil {
			newConn.Open(conn)
			return nil
		}

		// 报错
		newConn.Error(err)

		// 尝试重新连接
		if count > cfgCount {
			return err
		}

		// 等待
		time.Sleep(time.Duration(cfgInterval) * time.Second)
	}

	// return nil
}

// Stop close all listener
func (st *StreamTran) Stop() {
	st.running = false
	for _, listener := range st.listeners {
		listener.Close()
	}

	st.listeners = nil
}
