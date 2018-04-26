package soa

import (
	"fmt"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/codec"
	"github.com/jeckbjy/fairy/filter"
	"github.com/jeckbjy/fairy/frame"
	"github.com/jeckbjy/fairy/identity"
	"github.com/jeckbjy/fairy/tcp"
	"github.com/jeckbjy/fairy/util"
)

type zInfoGroup struct {
	infos []*Info // 信息组
	index int     // round索引
}
type zInfoMap map[uint64]*Info

// NewClient create client
func NewClient() *Client {
	cli := &Client{}
	return cli
}

/**
 * Client 注册或监听服务
 * 可以作为Provider,也可以作为Customer,或者两种同时存在
 * 作为Customer,需要注册一个Transport,因为通信协议需要外部指定
 * Customer连接成功后,可以调用GetInfoGroup,GetOneInfo,GetMinInfo,GetNextInfo
 *  函数查询订阅的服务器信息
 *
 * 用法:
 * step 1:注册信息
 *   11):Provider:PubAddr,PubPort,PubServices
 *   22):Customer:SubServices, SubTran
 * step 2:调用Start
 *
 */
type Client struct {
	id         uint64                 // 唯一ID
	centerTran fairy.Tran             // 注册中心Tran
	centerConn fairy.Conn             // 注册中心Conn
	info       *InfoEx                // 需要注册的信息
	subInfos   zInfoMap               // 订阅的服务Map
	subGroups  map[string]*zInfoGroup // 订阅的服务分组
	subTran    fairy.Tran             // 用于自动连接订阅的服务器
}

// PubAddr 注册Addr
func (cli *Client) PubAddr(outerAddr string, innerAddr string) {
	if cli.info == nil {
		cli.info = &InfoEx{}
	}

	if outerAddr != "" {
		cli.info.OuterAddr = outerAddr
	}

	if innerAddr != "" {
		cli.info.InnerAddr = innerAddr
	}
}

// PubPort 注册端口
func (cli *Client) PubPort(port uint) {
	if cli.info == nil {
		cli.info = &InfoEx{}
	}

	cli.info.Port = port
}

// PubServices 注册提供的服务
func (cli *Client) PubServices(services []string) {
	if cli.info == nil {
		cli.info = &InfoEx{}
	}

	cli.info.PubSerivces = append(cli.info.PubSerivces, services...)
}

// SubServices 注册订阅的服务
func (cli *Client) SubServices(services []string) {
	if cli.info == nil {
		cli.info = &InfoEx{}
	}

	cli.info.SubServices = append(cli.info.SubServices, services...)
}

// SubTran 设置Transport
func (cli *Client) SubTran(tran fairy.Tran) {
	tran.AddFilters(filter.NewConnect(cli.onConnectSub))
	cli.subTran = tran
}

// Start 连接服务器
func (cli *Client) Start(host string) error {
	// setup id
	id, err := util.NextID()
	if err != nil {
		return err
	}
	cli.id = id

	if cli.info == nil {
		return fmt.Errorf("must register info before start")
	}

	// setup inner ip
	if cli.info.InnerAddr == "" {
		addr, err := util.GetIPv4()
		if err != nil {
			return err
		}
		cli.info.InnerAddr = addr
	}

	// register callback
	register(&soaRegisterRsp{}, cli.onRegisterRsp)
	register(&soaRemoveMsg{}, cli.onRemoveMsg)

	// create transport
	tran := tcp.NewTran()
	tran.AddFilters(
		filter.NewFrame(frame.NewVarintLength()),
		filter.NewPacket(identity.NewInteger(), codec.NewGob()),
		filter.NewExecutor(),
		filter.NewConnect(cli.onConnect))

	tran.Connect(host, 0)
	tran.Start()

	cli.centerTran = tran
	return nil
}

// Stop 结束服务
func (cli *Client) Stop() {
	cli.centerTran.Stop()
}

func (cli *Client) addInfo(info *Info) {
	if cli.subGroups == nil {
		cli.subGroups = make(map[string]*zInfoGroup)
	}

	// add service to group
	for _, si := range info.PubSerivces {
		group, ok := cli.subGroups[si]
		if !ok {
			group = &zInfoGroup{}
			cli.subGroups[si] = group
		}

		group.infos = append(group.infos, info)
	}

	cli.subInfos[info.ID] = info
}

func (cli *Client) removeInfo(id uint64) {
	si, ok := cli.subInfos[id]
	if !ok {
		return
	}

	// remove
	for _, si := range si.PubSerivces {
		group, ok := cli.subGroups[si]
		if !ok {
			continue
		}

		for idx, info := range group.infos {
			if info.ID == id {
				group.infos = append(group.infos[:idx], group.infos[idx+1:]...)
				break
			}
		}
	}

	delete(cli.subInfos, id)
}

// 连接Center成功,自动注册自身信息
func (cli *Client) onConnect(conn fairy.Conn) {
	cli.centerConn = conn

	// 自动注册信息
	req := &soaRegisterReq{}
	req.Info = cli.info
	cli.centerConn.Send(req)
}

// 连接到订阅的服务器
func (cli *Client) onConnectSub(conn fairy.Conn) {
	// bind info
	id, ok := conn.GetTag().(uint64)
	if !ok {
		return
	}
	info, ok := cli.subInfos[id]
	if !ok {
		return
	}

	info.SetConn(conn)
}

func (cli *Client) onRegisterRsp(conn fairy.Conn, pkt fairy.Packet) {
	// 自动连接其他服务器
	rsp := pkt.GetMessage().(*soaRegisterRsp)
	// 注册
	for _, info := range rsp.SubInfos {
		// 不应该存在
		si, ok := cli.subInfos[info.ID]
		if ok {
			if si.GetConn() != nil {
				si.GetConn().Close()
			}

			*si = *info
		} else {
			cli.addInfo(info)
		}

		// 自动注册
		if cli.subTran != nil {
			cli.subTran.Connect(info.InnerAddr, info.ID)
		}
	}
}

func (cli *Client) onRemoveMsg(conn fairy.Conn, pkt fairy.Packet) {
	rsp := pkt.GetMessage().(*soaRemoveMsg)
	cli.removeInfo(rsp.ID)
}

func (cli *Client) onUpdateMsg(conn fairy.Conn, pkt fairy.Packet) {

}

// GetInfoGroup 获得Info Group
func (cli *Client) GetInfoGroup(name string) []*Info {
	group, ok := cli.subGroups[name]
	if ok {
		return group.infos
	}

	return nil
}

// GetOneInfo 查询第一个可用Info
func (cli *Client) GetOneInfo(name string) *Info {
	// 查询第一个
	group, ok := cli.subGroups[name]
	if ok && len(group.infos) > 0 {
		return group.infos[0]
	}

	return nil
}

// GetMinInfo 查询负载最小的Info
func (cli *Client) GetMinInfo(name string) *Info {
	group, ok := cli.subGroups[name]
	if !ok || len(group.infos) == 0 {
		return nil
	}

	min := group.infos[0]
	for i := 1; i < len(group.infos); i++ {
		if group.infos[i].Load < min.Load {
			min = group.infos[i]
		}
	}

	return min
}

// GetNextInfo 轮询下一个
func (cli *Client) GetNextInfo(name string) *Info {
	group, ok := cli.subGroups[name]
	if !ok || len(group.infos) == 0 {
		return nil
	}

	if group.index >= len(group.infos) {
		group.index = 0
	}

	idx := group.index
	group.index++
	return group.infos[idx]
}
