package itr

type Option func(server IServer)

type IServer interface {
	Start()                                                            // 启动服务器
	Stop()                                                             // 停止
	SetHandler(msgID uint32, handle IHandle)                           // 根据MsgID设置handle方法
	GetConnMgr() IConnMgr                                              // 获取server所有链接管理器
	GetPacket() IPacket                                                // 数据打包与解包对象
	SetPacket(IPacket)                                                 // 设置数据打包与解包对象
	SetType(svrType int32)                                             // 设置类型
	SetID(svrID int32)                                                 // 设置ID
	SetIP(ip string)                                                   // 设置IP
	SetPort(port int)                                                  // 设置端口
	SetMsgHandlerMgr(mgr IHandleMgr)                                   // 设置协议处理器
	SetConnectHookFunc(connected, disConnected func(conn IConnection)) // 设置网络连接方法
	CallConnStartHookFunc(conn IConnection)                            // 调用链接创建hook方法
	CallConnEndHookFunc(conn IConnection)                              // 调用链接断开hook方法
}

func WithSvrType(svrType int32) Option {
	return func(server IServer) {
		server.SetType(svrType)
	}
}

func WithSvrID(svrID int32) Option {
	return func(server IServer) {
		server.SetID(svrID)
	}
}

func WithIP(ip string) Option {
	return func(server IServer) {
		server.SetIP(ip)
	}
}

func WithPort(port int) Option {
	return func(server IServer) {
		server.SetPort(port)
	}
}

func WithPacket(packet IPacket) Option {
	return func(server IServer) {
		server.SetPacket(packet)
	}
}

func WithConnHookFunc(connected, disconnected func(conn IConnection)) Option {
	return func(server IServer) {
		server.SetConnectHookFunc(connected, disconnected)
	}
}

// server 基类实现
type BaseServer struct {
	svrType            int32
	sveID              int32
	ip                 string
	port               int
	onConnConnected    func(conn IConnection) // 简历链接hookFunc
	onConnDisconnected func(conn IConnection) // 断开链接hookFunc
	packet             IPacket                // 封包解包管理
	msgHandler         IHandleMgr             // 协议处理管理器
	connMgr            IConnMgr               // 链接管理器
}

var _ IServer = (*BaseServer)(nil)

func NewBaseServer() *BaseServer {
	return &BaseServer{}
}

func (b *BaseServer) Start() {
	panic("implement me")
}

func (b *BaseServer) Stop() {
	panic("implement me")
}

func (b *BaseServer) SetHandler(msgID uint32, handle IHandle) {
	if b.msgHandler == nil {
		return
	}

	b.msgHandler.RegisterHandle(msgID, handle)
}

func (b *BaseServer) GetConnMgr() IConnMgr {
	return b.connMgr
}

func (b *BaseServer) GetPacket() IPacket {
	return b.packet
}

func (b *BaseServer) SetPacket(packet IPacket) {
	if packet == nil {
		return
	}

	b.packet = packet
}

func (b *BaseServer) SetType(svrType int32) {
	b.svrType = svrType
}

func (b *BaseServer) SetID(svrID int32) {
	b.sveID = svrID
}

func (b *BaseServer) SetIP(ip string) {
	b.ip = ip
}

func (b *BaseServer) SetPort(port int) {
	b.port = port
}

func (b *BaseServer) SetMsgHandlerMgr(mgr IHandleMgr) {
	if mgr == nil {
		return
	}

	b.msgHandler = mgr
}

func (b *BaseServer) SetConnectHookFunc(connected, disConnected func(conn IConnection)) {
	b.onConnConnected = connected
	b.onConnDisconnected = disConnected
}

func (b *BaseServer) CallConnStartHookFunc(conn IConnection) {
	if b.onConnConnected == nil {
		return
	}

	b.onConnConnected(conn)
}

func (b *BaseServer) CallConnEndHookFunc(conn IConnection) {
	if b.onConnDisconnected == nil {
		return
	}

	b.onConnDisconnected(conn)
}
