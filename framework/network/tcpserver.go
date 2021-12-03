package network

import (
	"Doudou/framework/itr"
)

type TcpServer struct {
	*itr.BaseServer
}

var _ itr.IServer = (*TcpServer)(nil)

func NewTcpServer(ops ...itr.Option) itr.IServer {
	server := &TcpServer{
		BaseServer: itr.NewBaseServer(),
	}

	// server.SetPacket(NewNetPacket())
	// server.SetConnMgr(NewConnMgr())
	// server.SetMsgHandlerMgr(_default.NewApiMgr(1))

	for _, op := range ops {
		op(server)
	}

	return server
}

func (t *TcpServer) Start() {
	panic("implement me")
}

func (t *TcpServer) Stop() {
	panic("implement me")
}
