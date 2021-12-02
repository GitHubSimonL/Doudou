package framework

import (
	"Doudou/lib/logger"
	"fmt"
	"net"
)

type TCPServer struct {
	ServerBase
}

var _ IServer = (*TCPServer)(nil) // 编译期检查是否实现接口

func (ts *TCPServer) StartListen(port string, ops ...Option) {
	if ts.genSession == nil {
		logger.LogErrf("ServerReceiveConnect is nil", "")
		return
	}

	listenAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:"+port)
	if err != nil || listenAddr == nil {
		logger.LogErrf("start listen failed", fmt.Sprintf("err %v", err))
		return
	}

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil || listener == nil {
		logger.LogErrf("start listen failed", fmt.Sprintf("err %v", err))
		return
	}

	if ts.readMsgFunc == nil {
		logger.LogErrf("server readMsgFunc is nil", fmt.Sprintf("err %v", err))
		return
	}

	ts.SetState(SERVER_STATUS_RUNNING)
	go func() {
		for ts.GetState() == SERVER_STATUS_RUNNING {
			newConn, acErr := listener.Accept()
			if acErr != nil || newConn == nil {
				if newConn != nil {
					newConn.Close()
				}
				continue
			}

			if !ts.AccessCheck(newConn.RemoteAddr().String()) {
				logger.LogErrf("access denied", fmt.Sprintf("remote %v", newConn.RemoteAddr().String()))
				continue
			}

			newSession := ts.genSession(newConnection(newConn))
			if newSession == nil {
				newConn.Close()
				continue
			}

			newSession.SetSvrReceiveMsgChan(ts.GetReceiveMsgChan())
			newSession.SetReadMsgFunc(ts.readMsgFunc)
			newSession.Start()
		}
	}()

}

func NewTCPServerAgent(whiteListFile string, ops ...Option) (server *TCPServer) {
	newAgent := TCPServer{
		ServerBase: newServerBase(ops...),
	}

	if len(whiteListFile) > 0 && !newAgent.LoadWhiteList(whiteListFile) {
		logger.LogErrf("load whiteListFile failed", fmt.Sprintf("file %v", whiteListFile))
		return
	}

	return &newAgent
}
