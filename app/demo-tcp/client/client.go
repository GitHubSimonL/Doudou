package main

import (
	"Doudou/framework/itr"
	"Doudou/framework/network"
	_default "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func PingHandle(request itr.IRequest) {
	logger.LogDebugf("After Ping HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	request.GetConnection().SendMsg(2, request.GetData().([]byte))
}

func PongHandle(request itr.IRequest) {
	logger.LogDebugf("After Pong HandleMsg. Msg:%v Data:%v", request.GetMsgID(), request.GetData())
	request.GetConnection().SendMsg(1, []byte{byte(rand.Intn(100))})
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
	apiMgr.RegisterHandle(1, PingHandle)
	apiMgr.RegisterHandle(2, PongHandle)

	selfConn := network.NewConnection(nil, conn, 0, 1024, apiMgr, _default.NewNetPacket())
	go selfConn.Start()

	ts := time.NewTimer(2 * time.Second)
	defer ts.Stop()

	for {
		select {
		case <-ts.C:
			selfConn.SendMsg(1, []byte{byte(rand.Intn(100))})
		case <-selfConn.CloseSignal():
			return
		}
	}
}
