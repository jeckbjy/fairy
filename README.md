# Fairy library(WIP)

目标:高效,灵活,易用,易扩展的异步网络框架,设计上参考netty,mina,grizzly,使用责任链设计模式

- 支持tcp,websocket,kcp协议
- 支持protobuf,json,xml,gob编码
- 支持默认的消息处理线程模型
- 支持高效定时器

## 一:用例

```go
package chat

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/codecs"
	"github.com/jeckbjy/fairy/filters"
	"github.com/jeckbjy/fairy/frames"
	"github.com/jeckbjy/fairy/identities"
	"github.com/jeckbjy/fairy/log"
	"github.com/jeckbjy/fairy/tcp"
	"github.com/jeckbjy/fairy/timer"
	"github.com/jeckbjy/fairy/util"
)

type ChatMsg struct {
	Content   string
	Timestamp int64
}

func StartServer() {
	log.Debug("start server")
	// step1: register message
	fairy.RegisterMessage(&ChatMsg{}, nil)

	// step2: register handler
	fairy.RegisterHandler(&ChatMsg{}, func(ctx *fairy.HandlerCtx) {
		req := ctx.Message().(*ChatMsg)
		log.Debug("client msg:%+v", req)

		rsp := &ChatMsg{}
		rsp.Content = "welcome boy!"
		rsp.Timestamp = util.Now()
		ctx.Send(rsp)
	})

	// step3: create transport and add filters
	tran := tcp.NewTran()
	tran.AddFilters(
		filters.NewLogging(),
		filters.NewFrame(frames.NewLine()),
		filters.NewPacket(identities.NewString(), codecs.NewJson()),
		filters.NewExecutor())

	// step4: listen or connect
	tran.Listen(":8080", 0)
}

func StartClient() {
	log.Debug("start client")
	// step1: register message
	fairy.RegisterMessage(&ChatMsg{}, nil)

	// step2: register handler
	fairy.RegisterHandler(&ChatMsg{}, func(ctx *fairy.HandlerCtx) {
		req := ctx.Message().(*ChatMsg)
		log.Debug("server msg:%+v", req)
	})

	var gConn fairy.IConn
	// step3: create transport and add filters
	tran := tcp.NewTran()
	tran.AddFilters(
		filters.NewLogging(),
		filters.NewFrame(frames.NewLine()),
		filters.NewPacket(identities.NewString(), codecs.NewJson()),
		filters.NewExecutor())

	tran.AddFilters(filters.NewConnect(func(conn fairy.IConn) {
		// send msg to server
		req := &ChatMsg{}
		req.Content = "hello word!"
		conn.Send(req)
		gConn = conn
	}))

	// add timer for send message
	timer.Start(timer.ModeLoop, 1000, func() {
		log.Debug("Ontimeout")
		req := &ChatMsg{}
		req.Content = "hello word!"
		req.Timestamp = util.Now()
		gConn.Send(req)
	})

	// step4: listen or connect
	tran.Connect("localhost:8080", 0)
}
```

## 二:一些建议
- 服务器集群,对于复杂的服务器架构,直接使用默认的消息编码并不能满足需求,通常需要自定义IPacket和IIdentity来扩展,比如增加uid,消息源,目标类型等
- PacketFilter在某些情况下并不是高效的,因为里边进行了Codec的编解码,如果仅仅是转发协议,则并不需要解析body数据,可以自定义法Filter,通过判断是否有消息处理回调判断是否需要进行body解析
- rpc调用,本库并没有直接支持,如果需要,可以自定义Packet,增加一个唯一rpc id,Call时报错id到回调的映射,在消息处理处判断rpc id是否存在回调,如果存在则直接调用。额外需要一个定时器做延迟判断,防止消息永远没有返回,永远不能被执行
  
## 三:原理

- Transport和Connection
  - Transport:主要提供Listen和Connect两个接口,用于创建Connection,Connection默认会自动断线重连，如果不需要断线重连,可以通过SetOption关闭
  - Connection:类似于net.Conn，主要提供异步Read，Write，Close等接口

 ![Tran和Conn](doc/tran-conn.png)

- Filter
  - Filter 提供InBound和OutBound两种流向
    - InBound: HandleRead,HandleOpen,HandleError
    - OutBound:HandleWrite,HandleClose
  - FilterCtx 用于Filter之间数据传递,最常用的函数:GetData和SetData用于消息编解码,透传消息
  - 内置的filters
    - FrameFilter,PacketFilter,ExecutorFilter,LoggingFilter,TelnetFilter,ConnectFilter,RC4Filter
    - 自定义filter
      - filter应该是一个无状态的类,调用Next才会继续执行下一个,不调用将会终止传递
      - 如果需要数据，可以有两种方式：临时Filter之间传递数据，可以存储在FilterCtx中,长期持有的,可以存储在Connection中

![FilterChain](doc/filterchain.png)

- 消息的编解码
  - 在大部分应用中，消息的编解码是主要的通信工作，我这里划分了以下几个概念，Frame，Packet(Identity,Codec)
    - Frame:用于消息的粘包处理，例如类似http协议，以\r\n分隔，或者头部使用整数标识消息长度
    - Packet:消息包内容，通常分为两个部分，消息头和消息体,分别用Identity和Codec表示
      - Identity:用于消息头的编解码并创建具体的Packet
        - Fixed16Identity:小端编码,2个字节保存消息ID
        - StringIdentity:冒号分隔消息名和消息体
        - 自定义消息头:实现IIdentity接口并创建对应的IPacket
      - Codec:   用于消息体的编解码,例如json,protobuf

- 线程模型
  - Connection线程,每个Connection都会创建一个读和写协程
    - InBound在Connection的读线程中处理,直到转发到ExectorFilter逻辑线程中处理
    - Outbound在调用线程中处理,直到最终调用Write方法转到写协程中发送数据
  - 消息处理线程,并没有强制约定,可以自己继承Filter实现定制消息处理,默认发送到一个单独的消息处理协程中
    - 单线程模式:只需末尾添加ExectorFilter即可实现消息统一转发的Exector中的消息队列中执行
    - Executor可以不止一个线程,比如:某些复杂但又独立的业务操作，可以在注册消息回调时制定一个queueIndex,则可以实现该模块在独立的线程中执行，但要使用者自己保证线程安全
  - 其他线程:Log线程,Timer线程,Executor线程
    - log线程需要注意的是属性的初始化是非线程安全的，需要在主线程中设置属性，启动后将不能再修改
    - timer线程默认会将处理函数放到逻辑线程(Executor)中调用,如果不需要放到逻辑线程，可以将TimerEngine中的exector设置为nil

- 其他辅助类
  - buffer:底层的数据流存储，使用list存储[]byte，数据非连续的，可以像stream一样操作数据,使用时需要注意当前位置，以及哪些函数会影响当前位置
  - registry:非线程安全,用于消息的注册，可通过名字，或者id注册查询，也可以通过类型查询名字和id
  - dispatcher:非线程安全,handler的注册和查询

- 扩展:本项目不依赖任何库,均以插件的形式扩展
  - fairy-protobuf:protobuf扩展
    - https://github.com/jeckbjy/fairy-protobuf
    - https://github.com/golang/protobuf
  - fairy-kcp: kcp扩展
    - https://github.com/jeckbjy/fairy-kcp
    - https://github.com/xtaci/kcp-go
  - fairy-websocket: websocket扩展
    - https://github.com/jeckbjy/fairy-websocket
    - https://github.com/gorilla/websocket

- 参考框架
  - grizzly: https://javaee.github.io/grizzly/
    - 非常不错的java网络库框架
  - cellnet: https://github.com/davyxu/cellnet 
    - 功能很完善的go服务器框架，扩展性也还不错,支持websocket,tcp,kcp,以及各种编码协议,整体设计上还是很不错的服务器框架
  - leaf: https://github.com/name5566/leaf
    - 一个扩展性非常低go服务器框架,支持websocket,tcp以及protobuf，json编码，但架构设计上耦合太严重
