package kcp

import (
	"container/list"
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

func NewConn(tran fairy.Transport, filters fairy.FilterChain, side bool, kind int) *KcpConn {
	conn := &KcpConn{}
	conn.NewBase(tran, filters, side, kind)
	return conn
}

type KcpConn struct {
	base.Connection
	base.ConnReader
	base.ConnWriter
	conn net.Conn
	wg   sync.WaitGroup
}

func (kc *KcpConn) LocalAddr() net.Addr {
	return kc.conn.LocalAddr()
}

func (kc *KcpConn) RemoteAddr() net.Addr {
	return kc.conn.RemoteAddr()
}

func (kc *KcpConn) Open(conn net.Conn) {
	if kc.SwapState(fairy.ConnStateClosed, fairy.ConnStateOpen) {
		kc.conn = conn
		kc.NewWriter()
		kc.NewReader()
		kc.HandleOpen(kc)
		fairy.GetConnMgr().Put(kc)
		go kc.readThread()
	}
}

func (kc *KcpConn) Close() {
	// 线程安全调用
	if kc.SwapState(fairy.ConnStateOpen, fairy.ConnStateConnecting) {
		// 异步关闭，需要等待读写线程退出，才能退出
		go func() {
			kc.HandleClose(kc)
			kc.conn.Close()
			kc.StopWrite()
			kc.wg.Wait()
			kc.SetState(fairy.ConnStateClosed)
			kc.conn = nil
			// remove
			fairy.GetConnMgr().Remove(kc.GetConnId())
			// try reconnect
			trans := kc.GetTransport().(*KcpTran)
			trans.Reconnect(kc)
		}()
	}
}

func (kc *KcpConn) Read() *fairy.Buffer {
	return kc.Reader()
}

func (kc *KcpConn) Write(buffer *fairy.Buffer) {
	kc.PushBuffer(buffer, kc.sendThread)
}

func (kc *KcpConn) Send(msg interface{}) {
	kc.HandleWrite(kc, msg)
}

func (kc *KcpConn) readThread() {
	kc.wg.Add(1)
	defer kc.wg.Done()
	// loop read
	bufferSize := kc.GetConfig(fairy.KeyReaderBufferSize).(int)
	for {
		// 读取数据
		data := make([]byte, bufferSize)
		n, err := kc.conn.Read(data)
		if err == nil {
			kc.Append(data[:n])
			kc.HandleRead(kc)
		} else {
			kc.HandleError(kc, err)
			break
		}
	}
}

func (kc *KcpConn) sendThread() {
	kc.wg.Add(1)
	defer kc.wg.Done()

	for !kc.IsStopped() {
		buffers := list.List{}
		kc.WaitBuffers(&buffers)

		// write all buffer
		err := kc.WriteBuffers(kc.conn, &buffers)
		if err != nil {
			kc.HandleError(kc, err)
			break
		}
	}
}
