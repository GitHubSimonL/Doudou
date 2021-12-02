package network

import "Doudou/framework/itr"

type TcpServer struct {
	svrType            int32
	sveID              int32
	ip                 string
	port               int
	onConnConnected    func(conn itr.IConnection) // 简历链接hookFunc
	onConnDisconnected func(conn itr.IConnection) // 断开链接hookFunc
	packet             itr.IPacket                // 封包解包管理
	msgHandler         itr.IHandleMgr             // 协议处理管理器
	connMgr            itr.IConnMgr               // 链接管理器
}

var _ itr.IServer = (*TcpServer)(nil)

func NewTcpServer(ops ...itr.Option) itr.IServer {
	server := &TcpServer{
		packet: NewNetPacket(),
	}

	return server
}

func (t *TcpServer) Start() {
	panic("implement me")
}

func (t *TcpServer) Stop() {
	panic("implement me")
}

func (t *TcpServer) SetHandler(msgID int32, handle itr.IHandle) {
	panic("implement me")
}

func (t *TcpServer) GetConnMgr() itr.IConnMgr {
	panic("implement me")
}

func (t *TcpServer) SetConnStartHookFunc(f func(conn itr.IConnection)) {
	panic("implement me")
}

func (t *TcpServer) CallConnStartHookFunc(conn itr.IConnection) {
	panic("implement me")
}

func (t *TcpServer) SetConnEndHookFunc(f func(conn itr.IConnection)) {
	panic("implement me")
}

func (t *TcpServer) CallConnEndHookFunc(conn itr.IConnection) {
	panic("implement me")
}

func (t *TcpServer) Packet() itr.IPacket {
	panic("implement me")
}

func (t *TcpServer) SetIP(ip string) {
	t.ip = ip
}

func (t *TcpServer) SetPort(port int) {
	t.port = port
}

func (t *TcpServer) SetType(svrType int32) {
	t.svrType = svrType
}

func (t *TcpServer) SetID(svrID int32) {
	t.sveID = svrID
}
