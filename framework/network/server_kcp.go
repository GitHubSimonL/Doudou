package network

import (
	"Doudou/framework/itr"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"fmt"
	kcp "github.com/xtaci/kcp-go"
	"net"
)

type KcpServer struct {
	*itr.BaseServer
}

var _ itr.IServer = (*KcpServer)(nil)

func NewKcpServer(ops ...itr.Option) itr.IServer {
	server := &KcpServer{
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

func (u *KcpServer) Start() {
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

		kcpListener, err := kcp.ListenWithOptions(addr.String(), nil, 0, 0)
		if err != nil || kcpListener == nil {
			logger.LogWarnf("server start fail. err:%v", err.Error())
			return
		}

		var cid uint32 = 0
		for {
			conn, err := kcpListener.Accept()
			if err != nil || conn == nil {
				logger.LogWarnf("receive data failed", fmt.Sprintf("err %v", err))
				continue
			}

			if !u.AccessCheck(conn.RemoteAddr().String()) {
				logger.LogWarnf("blacklist ip connected. %v", conn.RemoteAddr().String())
				continue
			}

			// 最大链接保持数
			if u.GetConnMgr().Len() >= 100000 {
				logger.LogWarnf("The number of network links exceeds the threshold. %v", u.GetConnMgr().Len())
				conn.Close()
				continue
			}

			selfConn := NewConnection(u, conn, cid, 1024, u.GetApiMgr(), u.GetPacket())
			if selfConn == nil {
				conn.Close()
				continue
			}

			cid++
			go selfConn.Start()
		}
	}()
}

func (u *KcpServer) Stop() {
	defer func() {
		logger.LogDebugf("server stop finish.")
	}()

	u.GetConnMgr().ClearConn()

	u.StopSignal() <- struct{}{}
}
