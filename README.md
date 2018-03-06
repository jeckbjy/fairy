# Fairy library(WIP)

模仿netty,mina,grizzly，用go语言实现的一套网络库，api的设计上更接近于grizzly,为了方便游戏开发，除了网络模块，还有一些其他辅助功能，比如table,log,container等

## 架构
* Transport和Connection

    - Transport:主要提供Listen和Connect两个接口,用于创建Connection
    - Connection:类似于net.Conn，主要提供Read，Write，Close等接口，区别是读写是异步完成的

 ![Tran和Conn](doc/tran-conn.png)

 * FilterChain和Filters
    - 类似grizzly等，分为InBound和OutBound两种流向
    InBound: HandleRead,HandleOpen,HandleError
    OutBound:HandleWrite,HandleClose
 
    - 内置的filters
    FrameFilter,PacketFilter,ExecutorFilter,LoggingFilter,TelnetFilter,ConnectFilter
    - 自定义filter
    filter应该是一个无数据的类，如果需要数据，可以有两种方式：临时Filter之间传递数据，可以存储在FilterContext中，永久持有的，可以存储在Connection中

 * 消息的编解码
    - 在大部分应用中，消息的编解码是主要的通信工作，我这里划分了以下几个概念，Frame，Packet(Identity,Codec)
    - Frame:用于消息的粘包处理，例如类似http协议，以\r\n分隔，或者以固定长度的头标识后边消息长度
    - Packet:消息包内容，通常分为两个部分，消息头和消息体,分别用Identity和Codec表示
 Identity:最简单的两种形式,两个字节表示消息ID或者字符串表示消息名字
 Codec:消息具体的编解码,例如json,protobuf

 * 线程处理
    - 每个Connection都有自己的读写协程,
  InBound在Connection的读协程中处理,直到转发到ExectorFilter逻辑线程中处理
  Outbound在调用线程中处理,直到发送字节流时转到Connection的写协程中
    - 其他线程：Timer线程，Executor线程,Log线程等,
  log线程需要注意的是属性的初始化是非线程安全的，需要在主线程中设置属性，启动后将不能再修改
  timer线程默认会将处理函数放到逻辑线程(Executor)中调用,如果不需要放到逻辑线程，可以将TimerEngine中的exector设置为nil

* 其他辅助类
    - buffer:底层的数据流存储，使用list存储[]byte，数据时非连续的，可以像stream一样操作数据
    - registry:用于消息的注册，可通过名字，或者id注册查询，也可以通过类型查询名字和id
    - dispatcher:handler的注册和查询

* 辅助工具可参考
    - protobuf: https://github.com/jeckbjy/tool-proto-gen
    - 导表工具:  https://github.com/jeckbjy/tool-table-gen

* 依赖库
    - fairy:不依赖任何库
    - fairy-kcp:依赖 github.com/xtaci/kcp-go
    - fairy-protobuf: 依赖 github.com/golang/protobuf
    - fairy-websocket: 依赖 github.com/gorilla/websocket

* 参考框架
    - leaf: https://github.com/name5566/leaf 扩展性非常低
    - cellnet: https://github.com/davyxu/cellnet 
