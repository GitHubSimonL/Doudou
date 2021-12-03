package main

import (
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", _default.DefaultIP, _default.DefaultPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	selfConn := network.NewConnection(nil, conn, 0, 1024, _default.NewApiMgr(1), _default.NewNetPacket())
	go selfConn.Start()

	ts := time.NewTimer(1 * time.Second)
	idx := 0

	for {
		select {
		case <-ts.C:
			ts.Reset(1 * time.Second)
			selfConn.SendMsg(1, []byte{byte(idx)})

		}
	}

}
