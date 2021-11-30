package main

import (
	"Doudou/lib/logger"
	"Doudou/network"
	"fmt"
)

func init() {
	logger.InitLogger(
		logger.WithFilename("./log/demo1.log"),
		logger.WithPprof(true))
}

func main() {
	svr := network.NewTCPServerAgent("")
	svr.StartListen("11223")

	for {
		select {
		case msg, ok := <-svr.GetReceiveMsgChan():
			if !ok {
				continue
			}

			fmt.Println("=====", msg)

		}
	}
}
