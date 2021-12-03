package main

import (
	"Doudou/framework/itr"
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"fmt"
	"net"
	"time"
)

type Ping struct {
	itr.BaseHandle
}

func (p *Ping) AfterHandle(request itr.IRequest) {
	logger.LogDebugf("After Ping HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	request.GetConnection().SendMsg(2, request.GetData())
}

type Pong struct {
	itr.BaseHandle
}

func (p *Pong) AfterHandle(request itr.IRequest) {
	logger.LogDebugf("After Pong HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	request.GetConnection().SendMsg(1, request.GetData())
}

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

	apiMgr := _default.NewApiMgr(1)
	apiMgr.RegisterHandle(1, &Ping{})
	apiMgr.RegisterHandle(2, &Pong{})

	selfConn := network.NewConnection(nil, conn, 0, 1024, apiMgr, _default.NewNetPacket())
	go selfConn.Start()

	ts := time.NewTicker(5 * time.Second)
	defer ts.Stop()

	idx := 0

	for {
		select {
		case <-ts.C:
			selfConn.SendMsg(1, []byte{byte(idx)})
			idx++
		case <-selfConn.CloseSignal():
			return
		}
	}
}
