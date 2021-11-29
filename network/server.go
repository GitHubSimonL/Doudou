package network

import (
	"Doudou/lib/logger"
	"go.uber.org/atomic"
	"net"
	"runtime/debug"
)

type genSession func(conn net.Conn) ISession

// Server抽象接口
type IServer interface {
	Close()                                                                              // 关闭
	AfterClose(trigger func())                                                           // 关闭
	GetType() int32                                                                      // 类型
	GetID() int32                                                                        // ID
	LoadWhiteList(filename string) bool                                                  // 加载白名单
	AccessCheck(ip string) bool                                                          // 是否放行
	StartListen(port string, rDeadLine, wDeadLine int32, receiveFunc genSession) IServer // 开始监听
	SetState(state ServerState)
	GetState() ServerState
	GetReceiveMsgChan() chan NetMsg // 获取server接收消息channel，将它赋值给conn的的消息接收channel
}

type ServerState int32

const (
	SERVER_STATUS_LAUNCHING = ServerState(iota)
	SERVER_STATUS_RUNNING
	SERVER_STATUS_STOPPING
)

type ServerBase struct {
	svrType   int32
	svrID     int32
	closeChan chan struct{}
	state     *atomic.Int32
	WhiteList
	Sessions   map[int32]*BaseSession
	MsgChannel chan NetMsg
}

func (s *ServerBase) SetState(state ServerState) {
	s.state.Store(int32(state))
}

func NewServerBase() ServerBase {
	return ServerBase{
		closeChan: make(chan struct{}, 1),
		state:     atomic.NewInt32(int32(SERVER_STATUS_LAUNCHING)),
	}
}

func (s *ServerBase) AfterClose(trigger func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogErrf("panic recover %v:%v", r, string(debug.Stack()))
			}
		}()

		for {
			select {
			case _, ok := <-s.closeChan:
				if !ok {
					continue
				}

				if trigger != nil {
					trigger()
				}

				return
			}
		}
	}()
}

func (s *ServerBase) Close() {
	s.closeChan <- struct{}{}
}

func (s *ServerBase) GetType() int32 {
	return s.svrType
}

func (s *ServerBase) GetState() ServerState {
	return ServerState(s.state.Load())
}

func (s *ServerBase) GetID() int32 {
	return s.svrID
}

func (s *ServerBase) GetReceiveMsgChan() chan NetMsg {
	return s.MsgChannel
}
