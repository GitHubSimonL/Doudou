package network

import (
	. "Doudou/framework/itr"
)

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

func WithConnMgr(mgr IConnMgr) Option {
	return func(server IServer) {
		server.SetConnMgr(mgr)
	}
}

func WithApiMgr(mgr IApiMgr) Option {
	return func(server IServer) {
		server.SetMsgHandlerMgr(mgr)
	}
}
