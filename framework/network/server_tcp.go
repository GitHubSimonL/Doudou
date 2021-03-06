package network

import (
	"Doudou/framework/itr"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"fmt"
	"net"
)

type TcpServer struct {
	*itr.BaseServer
}

var _ itr.IServer = (*TcpServer)(nil)

func NewTcpServer(ops ...itr.Option) itr.IServer {
	server := &TcpServer{
		BaseServer: itr.NewBaseServer(),
	}

	server.SetPort(_default.DefaultPort)
	server.SetIP(_default.DefaultIP)
	server.SetPacket(_default.NewNetPacket())
	server.SetConnMgr(_default.NewConnMgr())
	server.SetMsgHandlerMgr(_default.NewApiMgr(1))

	for _, op := range ops {
		op(server)
	}

	return server
}

func (t *TcpServer) Start() {
	defer func() {
		logger.LogDebugf("server start finish. ip:%v port:%v", t.GetIP(), t.GetPort())
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogErrf("server work catch err. %v", err)
			}
		}()

		t.GetApiMgr().StartWorkPool()

		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", t.GetIP(), t.GetPort()))
		if err != nil {
			logger.LogWarnf("server start fail. err:%v", err.Error())
			return
		}

		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			logger.LogWarnf("server start fail. err:%v", err.Error())
			return
		}

		var cid uint32 = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logger.LogWarnf("tcp listener catch err. %v", err)
				continue
			}

			// 最大链接保持数
			if t.GetConnMgr().Len() >= 100000 {
				logger.LogWarnf("The number of network links exceeds the threshold. %v", t.GetConnMgr().Len())
				conn.Close()
				continue
			}

			selfConn := NewConnection(t, conn, cid, 1024, t.GetApiMgr(), t.GetPacket())
			if selfConn == nil {
				conn.Close()
				continue
			}

			if !t.AccessCheck(selfConn.RemoteAddr().String()) {
				logger.LogWarnf("blacklist ip connected. %v", selfConn.RemoteAddr().String())
				conn.Close()
				continue
			}

			cid++
			go selfConn.Start()
		}
	}()
}

func (t *TcpServer) Stop() {
	defer func() {
		logger.LogDebugf("server stop finish.")
	}()

	t.GetConnMgr().ClearConn()

	t.StopSignal() <- struct{}{}
}
