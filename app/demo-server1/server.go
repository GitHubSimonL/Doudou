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

func (p *Ping) Handle(request itr.IRequest) {
	logger.LogDebugf("After Ping HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	time.Sleep(1 * time.Second)
	request.GetConnection().SendMsg(2, request.GetData())
}

type Pong struct {
	itr.BaseHandle
}

func (p *Pong) Handle(request itr.IRequest) {
	logger.LogDebugf("After Pong HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	time.Sleep(1 * time.Second)
	request.GetConnection().SendMsg(1, request.GetData())
}

func main() {
	localServer := network.NewTcpServer(
		network.WithConnMgr(_default.NewConnMgr()),
		network.WithApiMgr(_default.NewApiMgr(1)),
		network.WithPacket(_default.NewNetPacket()),
	)

	localServer.Start()
	localServer.SetHandler(1, &Ping{})
	localServer.SetHandler(2, &Pong{})

	go func() {
		time.Sleep(30 * time.Minute)
		localServer.Stop()
	}()

	for {
		select {
		case <-localServer.StopSignal():
			return
		case req, ok := <-localServer.ReadReq():
			if !ok {
				return
			}

			if req == nil {
				continue
			}

			localServer.GetApiMgr().AddMgsToTaskPool(req)
		}
	}
}
