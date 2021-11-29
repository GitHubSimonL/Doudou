package network

import (
	"Doudou/lib/logger"
	"fmt"
	"net"
)

type TCPServer struct {
	ServerBase
}

func (ts *TCPServer) StartListen(port string, receiveFunc genSession) {
	if receiveFunc == nil {
		logger.LogErrf("ServerReceiveConnect is nil", "")
		return
	}

	listenAddr, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil || listenAddr == nil {
		logger.LogErrf("start listen failed", fmt.Sprintf("err %v", err))
		return
	}

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil || listener == nil {
		logger.LogErrf("start listen failed", fmt.Sprintf("err %v", err))
		return
	}

	ts.SetState(SERVER_STATUS_RUNNING)
	go func() {
		for ts.GetState() == SERVER_STATUS_RUNNING {
			newConn, acErr := listener.Accept()
			if acErr != nil || newConn == nil {
				continue
			}

			if !ts.AccessCheck(newConn.RemoteAddr().String()) {
				logger.LogErrf("access denied", fmt.Sprintf("remote %v", newConn.RemoteAddr().String()))
				continue
			}

			newSession := receiveFunc(newConn)
			if newSession == nil {
				continue
			}

			newSession.SetSvrReceiveMsgChan(ts.GetReceiveMsgChan())
			newSession.Start()
		}
	}()

}

func NewTCPServerAgent(whiteListFile string) (server *TCPServer) {
	newAgent := TCPServer{
		ServerBase: NewServerBase(),
	}

	if len(whiteListFile) > 0 && !newAgent.LoadWhiteList(whiteListFile) {
		logger.LogErrf("load whiteListFile failed", fmt.Sprintf("file %v", whiteListFile))
		return
	}

	return &newAgent
}
