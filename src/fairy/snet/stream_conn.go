package snet

import (
	"container/list"
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

func NewConn(channel IChannel, tran fairy.Transport, side bool, kind int) *StreamConn {
	stream_conn := &StreamConn{}
	stream_conn.Create(channel, tran, side, kind)
	return stream_conn
}

type IChannel interface {
	Read(cap int) ([]byte, error)
	Write(buf []byte) error
	Open(conn interface{})
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

type StreamConn struct {
	base.Conn
	channel  IChannel
	wg       sync.WaitGroup
	rbuf     *fairy.Buffer
	wbuf     *list.List
	wcnd     *sync.Cond
	wmux     *sync.Mutex
	wstopped bool
}

func (sc *StreamConn) Create(channel IChannel, tran fairy.Transport, side bool, kind int) {
	sc.Conn.Create(tran, side, kind)
	sc.channel = channel
	sc.rbuf = fairy.NewBuffer()
	sc.wmux = &sync.Mutex{}
	sc.wstopped = true
}

func (sc *StreamConn) LocalAddr() net.Addr {
	return sc.channel.LocalAddr()
}

func (sc *StreamConn) RemoteAddr() net.Addr {
	return sc.channel.RemoteAddr()
}

func (sc *StreamConn) Error(err error) {
	sc.HandleError(sc, err)
}

func (sc *StreamConn) Open(conn interface{}) {
	if !sc.SwapState(fairy.ConnStateClosed, fairy.ConnStateOpen) {
		return
	}

	sc.channel.Open(conn)
	// set reader
	sc.rbuf.Clear()

	sc.HandleOpen(sc)
	fairy.GetConnMgr().Put(sc)
	go sc.recvThread()
}

func (sc *StreamConn) Close() {
	if !sc.SwapState(fairy.ConnStateOpen, fairy.ConnStateConnecting) {
		return
	}

	go func() {
		sc.channel.Close()
		// stop write, wait for finish
		sc.wstopped = true
		sc.wcnd.Signal()
		sc.wg.Wait()
		sc.SetState(fairy.ConnStateClosed)
		//
		fairy.GetConnMgr().Remove(sc.GetConnId())
		// reconnect
		if sc.IsClientSide() {
			sc.GetTransport().Reconnect(sc)
		}
	}()
}

func (sc *StreamConn) Wait() {
	sc.wg.Wait()
}

func (sc *StreamConn) Read() *fairy.Buffer {
	return sc.rbuf
}

func (sc *StreamConn) Write(buf *fairy.Buffer) {
	sc.wmux.Lock()

	// lazy init writer buffer
	if sc.wbuf == nil {
		sc.wbuf = list.New()
		sc.wcnd = sync.NewCond(sc.wmux)
	}

	if sc.wstopped {
		go sc.sendThread()
	}

	sc.wbuf.PushBack(buf)
	// sc.wfuture
	sc.wcnd.Signal()
	sc.wmux.Unlock()
}

func (sc *StreamConn) Send(msg interface{}) {
	sc.HandleWrite(sc, msg)
}

func (sc *StreamConn) recvThread() {
	// log.Debug("recv thread start")

	sc.wg.Add(1)

	bufSize := sc.GetConfig(fairy.CfgReaderBufferSize).(int)
	for {
		data, err := sc.channel.Read(bufSize)
		if err != nil {
			sc.HandleError(sc, err)
			break
		}

		sc.rbuf.Append(data)
		sc.HandleRead(sc)
	}

	sc.wg.Done()
	// log.Debug("recv thread finish")
}

func (sc *StreamConn) sendThread() {
	// log.Debug("send thread start")

	sc.wg.Add(1)

	sc.wstopped = false
	for !sc.wstopped {
		bufs := list.List{}
		// wait buffer
		sc.wmux.Lock()
		for !sc.wstopped && sc.wbuf.Len() == 0 {
			sc.wcnd.Wait()
		}

		bufs = *sc.wbuf
		sc.wbuf.Init()
		sc.wmux.Unlock()

		// flush buffer
		for iterl := bufs.Front(); iterl != nil; iterl = iterl.Next() {
			if iterl.Value == nil {
				continue
			}

			buffer := iterl.Value.(*fairy.Buffer)
			iterb := buffer.Front()
			for ; iterb != nil; iterb = iterb.Next() {
				data := iterb.Value.([]byte)
				err := sc.channel.Write(data)
				if err != nil {
					sc.HandleError(sc, err)
					sc.wstopped = true
					break
				}
			}

			if iterb != nil {
				break
			}
		}
	}

	sc.wbuf.Init()
	sc.wg.Done()
	// log.Debug("send thread finish")
}
