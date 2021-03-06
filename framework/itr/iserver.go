package itr

import "Doudou/lib/logger"

type Option func(server IServer)

type IServer interface {
	Start()                                                            // 启动服务器
	Stop()                                                             // 停止
	SetHandler(msgID uint32, handle HandleFunc)                        // 根据MsgID设置handle方法
	GetApiMgr() IApiMgr                                                // 协议处理管理器
	GetPacket() IPacket                                                // 数据打包与解包对象
	SetPacket(IPacket)                                                 // 设置数据打包与解包对象
	SetMsgHandlerMgr(mgr IApiMgr)                                      // 设置协议处理器
	GetConnMgr() IConnMgr                                              // 获取server所有链接管理器
	SetConnMgr(mgr IConnMgr)                                           // 设置链接管理器
	SetConnectHookFunc(connected, disConnected func(conn IConnection)) // 设置网络连接方法
	CallConnStartHookFunc(conn IConnection)                            // 调用链接创建hook方法
	CallConnEndHookFunc(conn IConnection)                              // 调用链接断开hook方法
	LoadWhiteList(filename string) bool                                // 加载白名单
	AccessCheck(ip string) bool                                        // 是否放行
	ReadReq() chan IRequest                                            // 读客户端请求
	WriteReq(req IRequest)

	SetType(svrType int32) // 设置类型
	SetID(svrID int32)     // 设置ID
	SetIP(ip string)       // 设置IP
	SetPort(port int)      // 设置端口
	GetType() int32        // 获取类型
	GetID() int32          // 获取ID
	GetIP() string         // 获取IP
	GetPort() int          // 获取端口

	StopSignal() chan struct{}
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
	apiMgr             IApiMgr                // 协议处理管理器
	connMgr            IConnMgr               // 链接管理器
	stopSignal         chan struct{}
	reqList            chan IRequest
	WhiteList
}

func (b *BaseServer) WriteReq(req IRequest) {
	if req == nil || len(b.reqList) >= 256 {
		logger.LogWarnf("req is nil or list too long. %v", len(b.reqList))
	}

	b.reqList <- req
}

var _ IServer = (*BaseServer)(nil)

func NewBaseServer() *BaseServer {
	return &BaseServer{
		stopSignal: make(chan struct{}, 1),
		reqList:    make(chan IRequest, 256),
	}
}

func (b *BaseServer) Start() {
	panic("implement me")
}

func (b *BaseServer) Stop() {
	panic("implement me")
}

func (b *BaseServer) SetHandler(msgID uint32, handle HandleFunc) {
	if b.apiMgr == nil {
		return
	}

	b.apiMgr.RegisterHandle(msgID, handle)
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

func (b *BaseServer) SetMsgHandlerMgr(mgr IApiMgr) {
	if mgr == nil {
		return
	}

	b.apiMgr = mgr
}

func (b *BaseServer) SetConnMgr(mgr IConnMgr) {
	b.connMgr = mgr
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

func (b *BaseServer) GetApiMgr() IApiMgr {
	return b.apiMgr
}

func (b *BaseServer) GetType() int32 {
	return b.svrType
}

func (b *BaseServer) GetID() int32 {
	return b.sveID
}

func (b *BaseServer) GetIP() string {
	return b.ip
}

func (b *BaseServer) GetPort() int {
	return b.port
}

func (b *BaseServer) StopSignal() chan struct{} {
	return b.stopSignal
}

func (b *BaseServer) ReadReq() chan IRequest {
	return b.reqList
}
