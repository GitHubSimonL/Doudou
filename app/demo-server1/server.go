package main

import (
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
)

func main() {
	servr := network.NewTcpServer(
		network.WithConnMgr(_default.NewConnMgr()),
		network.WithApiMgr(_default.NewApiMgr(1)),
		network.WithPacket(_default.NewNetPacket()),
	)

	servr.Start()

	select {}

}
