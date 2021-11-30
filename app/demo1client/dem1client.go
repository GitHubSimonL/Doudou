package main

import (
	"net"
	"time"
)

func main() {
	listenAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:"+"11223")
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, listenAddr)
	if err != nil {
		return
	}

	for {
		time.Sleep(1 * time.Second)
		conn.Write([]byte{1, 2})
	}
}
