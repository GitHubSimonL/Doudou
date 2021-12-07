package network

import (
	"Doudou/framework/itr"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"fmt"
	"net"
)

type UdpServer struct {
	*itr.BaseServer
}

var _ itr.IServer = (*UdpServer)(nil)

func NewUdpServer(ops ...itr.Option) itr.IServer {
	server := &UdpServer{
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

func (u *UdpServer) Start() {
	defer func() {
		logger.LogDebugf("server start finish. ip:%v port:%v", u.GetIP(), u.GetPort())
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogErrf("server work catch err. %v", err)
			}
		}()

		u.GetApiMgr().StartWorkPool()

		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", u.GetIP(), u.GetPort()))
		if err != nil {
			logger.LogWarnf("server start fail. err:%v", err.Error())
			return
		}

		udpConn, err := net.ListenUDP("udp", addr)
		if err != nil || udpConn == nil {
			logger.LogWarnf("server start fail. err:%v", err.Error())
			return
		}
		
		var cid uint32 = 0
		for {
			conn, err := listener.()
			if err != nil {
				logger.LogWarnf("tcp listener catch err. %v", err)
				continue
			}

			// 最大链接保持数
			if u.GetConnMgr().Len() >= 100000 {
				logger.LogWarnf("The number of network links exceeds the threshold. %v", u.GetConnMgr().Len())
				conn.Close()
				continue
			}

			selfConn := NewConnection(t, conn, cid, 1024, u.GetApiMgr(), u.GetPacket())
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

func (u *UdpServer) Stop() {
	defer func() {
		logger.LogDebugf("server stop finish.")
	}()

	u.GetConnMgr().ClearConn()

	u.StopSignal() <- struct{}{}
}
