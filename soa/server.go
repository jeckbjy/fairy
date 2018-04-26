package soa

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/codec"
	"github.com/jeckbjy/fairy/filter"
	"github.com/jeckbjy/fairy/frame"
	"github.com/jeckbjy/fairy/identity"
	"github.com/jeckbjy/fairy/tcp"
	"github.com/jeckbjy/fairy/util"
)

// server内使用,用于快速比较订阅的服务
// 每个key会映射到一个位索引
type srvInfo struct {
	InfoEx
	pubs []byte
	subs []byte
}

// 自己是否订阅了他人的服务,求交集
func (info *srvInfo) isSubscribe(other *srvInfo) bool {
	num := util.MinInt(len(info.subs), len(other.pubs))
	for i := 0; i < num; i++ {
		if (info.subs[i] & other.pubs[i]) != 0 {
			return true
		}
	}

	return false
}

// NewServer create server
func NewServer() *Server {
	srv := &Server{}
	srv.infos = make(map[uint64]*srvInfo)
	srv.keyMap = make(map[string]uint)
	return srv
}

// Server 负责节点的注册管理以及分发
/**
 * 简单的服务治理,cs结构
 * server管理client服务信息,当有新的链接,或链接断开时,通知相关的client
 * client注册自己提供的服务，以及订阅的服务,并自动连接订阅的服务器
 * client作为Provider时,目前仅可以有一个端口给其他Customer连接
 */
type Server struct {
	tran   fairy.Tran
	infos  map[uint64]*srvInfo
	keyMap map[string]uint
	keyID  uint
}

// Start 启动服务
func (srv *Server) Start(host string) {
	// register callback
	register(soaRegisterReq{}, srv.onRegister)

	// create transport
	tran := tcp.NewTran()
	tran.AddFilters(
		filter.NewFrame(frame.NewVarintLength()),
		filter.NewPacket(identity.NewInteger(), codec.NewGob()),
		filter.NewExecutor(),
		filter.NewClose(srv.handleClose))

	tran.Listen(host, 0)
	tran.Start()

	srv.tran = tran
}

// Stop 结束服务
func (srv *Server) Stop() {
	srv.tran.Stop()
}

func (srv *Server) remove(info *srvInfo) {
	// 通知删除服务
	// broadcast
	rsp := &soaRemoveMsg{}
	rsp.ID = info.ID
	for _, si := range srv.infos {
		if si.isSubscribe(info) {
			si.GetConn().Send(rsp)
		}
	}
}

func (srv *Server) handleClose(conn fairy.Conn) {
	// 通知删除服务器
	infoID, ok := conn.GetData().(uint64)
	if !ok {
		return
	}

	// find and remove
	if info, ok := srv.infos[infoID]; ok {
		srv.remove(info)
	}
}

func (srv *Server) setKeyMap(keyStr []string, keyBit *[]byte) {
	for _, key := range keyStr {
		id, ok := srv.keyMap[key]
		if !ok {
			id = srv.keyID
			srv.keyID++
			srv.keyMap[key] = id
		}

		// set bits
		num := id / 8
		bit := id % 8
		if num > uint(len(*keyBit)) {
			newBit := make([]byte, num+1)
			copy(newBit, *keyBit)
			*keyBit = newBit
		}

		(*keyBit)[num] |= byte(1 << bit)
	}
}

// 新的客户端注册
func (srv *Server) onRegister(conn fairy.Conn, pkt fairy.Packet) {
	req := pkt.GetMessage().(soaRegisterReq)
	info := &srvInfo{}
	info.InfoEx = *req.Info
	srv.setKeyMap(info.PubSerivces, &info.pubs)
	srv.setKeyMap(info.SubServices, &info.subs)

	// try remove old
	if oldInfo, ok := srv.infos[info.ID]; ok {
		srv.remove(oldInfo)
	}

	srv.infos[info.ID] = info
	conn.SetData(info.ID)

	// 通知自己关注的服务
	if len(info.SubServices) != 0 {
		rsp := &soaRegisterRsp{}
		for _, sinfo := range srv.infos {
			//
			if info.isSubscribe(sinfo) {
				rsp.SubInfos = append(rsp.SubInfos, &sinfo.Info)
			}
		}
	}

	// 通知其他关注自己的服务
	if len(info.PubSerivces) != 0 {
		rsp := &soaRegisterRsp{}
		rsp.SubInfos = append(rsp.SubInfos, &info.Info)
		for _, sinfo := range srv.infos {
			//
			if sinfo.isSubscribe(info) {
				sinfo.GetConn().Send(rsp)
			}
		}
	}
}

// 更新数据
func (srv *Server) onUpdate(conn fairy.Conn, pkt fairy.Packet) {
	// 更新数据,并通知其他人
}
