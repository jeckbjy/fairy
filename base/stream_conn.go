package base

import (
	"container/list"
	"net"
	"sync"
	"sync/atomic"

	"github.com/jeckbjy/fairy"
)

const (
	connStateOpened  = 0 // 已连接
	connStateClosed  = 1 // 已关闭
	connStateClosing = 2 // 关闭中,主动关闭,等待写操作完成
)

// NewConn 创建StreamConn
func NewConn(tran fairy.ITran, connector bool, tag string) *StreamConn {
	sc := &StreamConn{}
	sc.Conn.Init(tran, connector, tag)
	// init
	sc.state = connStateClosed
	sc.rbuf = fairy.NewBuffer()
	sc.wbuf = list.New()
	sc.wcnd = sync.NewCond(&sc.mutex)
	return sc
}

// StreamConn 长连接Connection基类
type StreamConn struct {
	Conn
	channel net.Conn
	wg      sync.WaitGroup
	rbuf    *fairy.Buffer // 读缓存
	wbuf    *list.List    // 写缓存
	wcnd    *sync.Cond    // 写同步
	wstop   bool          // 写线程是否存在
	rstop   bool          // 读线程是否存在
	state   int32         // 状态
	mutex   sync.Mutex    // 保护state,channel和wbuf
}

func (sc *StreamConn) IsActive() bool {
	return atomic.LoadInt32(&sc.state) == connStateOpened
}

func (sc *StreamConn) LocalAddr() net.Addr {
	return sc.channel.LocalAddr()
}

func (sc *StreamConn) RemoteAddr() net.Addr {
	return sc.channel.RemoteAddr()
}

func (sc *StreamConn) Open(conn net.Conn) {
	if conn == nil {
		return
	}

	closed := false

	sc.mutex.Lock()

	if sc.state == connStateClosed {
		closed = true
		sc.state = connStateOpened
		sc.channel = conn
		sc.rbuf.Clear()
		sc.wbuf.Init()
		go sc.recvThread()
		go sc.sendThread()
	}

	sc.mutex.Unlock()

	if closed {
		// unlock之后调用,防止回调里调用Close导致死锁
		sc.GetChain().HandleOpen(sc)
	}
}

// Close 主动关闭
func (sc *StreamConn) Close() {
	sc.mutex.Lock()

	if sc.state == connStateOpened {
		if !sc.wstop {
			// 等待写关闭
			sc.state = connStateClosing
			sc.wstop = true
			sc.wcnd.Signal()
		} else {
			// 直接关闭
			sc.state = connStateClosed
			if sc.channel != nil {
				sc.channel.Close()
				sc.channel = nil
			}
		}
	}

	sc.mutex.Unlock()
}

// 异常关闭,会尝试断线重连
func (sc *StreamConn) doClose(err error) {
	reconnect := false

	sc.mutex.Lock()

	// 非主动关闭的connector才尝试断线重连
	if sc.state == connStateOpened && sc.IsConnector() && err != nil {
		reconnect = true
	}

	// 可能是Open状态也可能是Closing状态
	if sc.state != connStateClosed {
		sc.state = connStateClosed

		if !sc.rstop {
			// 通知关闭读线程
			sc.channel.Close()
		}

		if !sc.wstop {
			// 通知关闭写线程
			sc.wstop = true
			sc.wcnd.Signal()
		}

		sc.channel = nil
	}

	sc.mutex.Unlock()

	if err != nil {
		// io.EOF其实是正常关闭socket
		sc.Error(err)
	}

	// 尝试断线重连
	if reconnect {
		sc.GetTran().(*StreamTran).Reconnect(sc)
	}
}

func (sc *StreamConn) Read() *fairy.Buffer {
	return sc.rbuf
}

func (sc *StreamConn) Write(buf *fairy.Buffer) error {
	sc.mutex.Lock()
	if sc.state == connStateOpened {
		sc.wbuf.PushBack(buf)
		sc.wcnd.Signal()
	}
	sc.mutex.Unlock()

	return nil
}

func (sc *StreamConn) Send(msg interface{}) error {
	sc.GetChain().HandleWrite(sc, msg)
	return nil
}

// Error 发生错误
func (sc *StreamConn) Error(err error) {
	sc.GetChain().HandleError(sc, err)
}

// 异步读线程
func (sc *StreamConn) recvThread() {
	sc.wg.Add(1)
	defer sc.wg.Done()

	sc.rstop = false
	// bufSize := 1024
	var err error
	for {
		err = sc.rbuf.ReadAll(sc.channel)
		if err != nil {
			break
		}

		sc.GetChain().HandleRead(sc)
	}
	sc.rstop = true

	// 关闭并自动尝试重新连接
	sc.doClose(err)
}

// 异步写线程
func (sc *StreamConn) sendThread() {
	sc.wg.Add(1)
	defer sc.wg.Done()

	sc.wstop = false
	var err error
	for {
		buffers := list.List{}
		sc.mutex.Lock()
		for !sc.wstop && sc.wbuf.Len() == 0 {
			sc.wcnd.Wait()
		}

		buffers = *sc.wbuf
		sc.wbuf.Init()
		sc.mutex.Unlock()

		if sc.wstop {
			break
		}

		// flush
		for iter := buffers.Front(); iter != nil; iter = iter.Next() {
			if iter.Value == nil {
				continue
			}

			buffer := iter.Value.(*fairy.Buffer)
			err = buffer.WriteAll(sc.channel)
			if err != nil {
				break
			}
		}

		if err != nil {
			break
		}
	}

	sc.wstop = true
	sc.doClose(err)
}
