package main

import (
	"Doudou/framework/itr"
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"time"
)

func PingHandle(request itr.IRequest) {
	logger.LogDebugf("Ping HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())

	time.Sleep(1 * time.Second)
	request.GetConnection().SendMsg(2, request.GetData().([]byte))
}

func PongHandle(request itr.IRequest) {
	logger.LogDebugf("Pong HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())

	time.Sleep(1 * time.Second)
	request.GetConnection().SendMsg(1, request.GetData().([]byte))
}

func main() {
	localServer := network.NewUdpServer(
		network.WithConnMgr(_default.NewConnMgr()),
		network.WithApiMgr(_default.NewApiMgr(1)),
		network.WithPacket(_default.NewNetPacket()),
	)

	localServer.Start()
	localServer.SetHandler(1, PingHandle)
	localServer.SetHandler(2, PongHandle)

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
