package main

import (
	"Doudou/framework/itr"
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"time"
)

type Ping struct {
	itr.BaseHandle
}

func (p *Ping) AfterHandle(request itr.IRequest) {
	logger.LogDebugf("After Ping HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	time.Sleep(1 * time.Second)
	request.GetConnection().SendMsg(2, request.GetData())
}

type Pong struct {
	itr.BaseHandle
}

func (p *Pong) AfterHandle(request itr.IRequest) {
	logger.LogDebugf("After Pong HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	time.Sleep(1 * time.Second)
	request.GetConnection().SendMsg(1, request.GetData())
}

func main() {
	servr := network.NewTcpServer(
		network.WithConnMgr(_default.NewConnMgr()),
		network.WithApiMgr(_default.NewApiMgr(1)),
		network.WithPacket(_default.NewNetPacket()),
	)

	servr.Start()
	servr.SetHandler(1, &Ping{})
	servr.SetHandler(2, &Pong{})

	go func() {
		time.Sleep(30 * time.Minute)
		servr.Stop()
	}()

	select {
	case <-servr.StopSignal():
		return
	}
}
