package main

import (
	"Doudou/framework/itr"
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
	"fmt"
)

type TestH struct {
	itr.BaseHandle
}

func (th *TestH) Handle(request itr.IRequest) {
	fmt.Printf("ConnID:%v MsgID:%v Data:%v", request.GetConnection().GetConnID(), request.GetMsgID(), request.GetData())
}

func main() {
	servr := network.NewTcpServer(
		network.WithConnMgr(_default.NewConnMgr()),
		network.WithApiMgr(_default.NewApiMgr(1)),
		network.WithPacket(_default.NewNetPacket()),
	)

	servr.Start()
	servr.SetHandler(1, &TestH{})
	select {}
}
