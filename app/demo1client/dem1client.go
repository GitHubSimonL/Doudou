package main

import (
	"bufio"
	"fmt"
	"io"
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

	exitChan := make(chan bool, 1)
	idx := 0

	rd := bufio.NewReader(conn)
	go func() {
		for {
			headerBuff := make([]byte, 13)
			_, err := io.ReadFull(rd, headerBuff)

			if err != nil {
				fmt.Printf("isEnd:%v conn err. %v \n", err == io.EOF, err)
				exitChan <- true
			}
		}
	}()

	sendTimer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-sendTimer.C:
			// _, err := conn.Write([]byte{1, 2})
			// if err != nil {
			// 	fmt.Printf("wirte err:%v \n", err)
			// 	return
			// }
			idx++
			fmt.Println("send: ", idx)

		case <-exitChan:
			return
		}
	}
}
