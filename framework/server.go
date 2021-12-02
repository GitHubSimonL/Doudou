package framework

import (
	"Doudou/lib/logger"
	"go.uber.org/atomic"
	"io"
	"runtime/debug"
	"time"
)

type genSession func(conn ICon) ISession
type readMsgFunc func(rd io.Reader) INetMsg

// Server抽象接口
type IServer interface {
	Close()                                     // 关闭
	AfterClose(callback func())                 // 关闭后回调，注意（这里的callback是线程不安全的）
	GetType() int32                             // 类型
	GetID() int32                               // ID
	LoadWhiteList(filename string) bool         // 加载白名单
	AccessCheck(ip string) bool                 // 是否放行
	StartListen(port string, options ...Option) // 开始监听
	SetState(state ServerState)
	GetState() ServerState
	GetReceiveMsgChan() chan INetMsg // 获取server接收消息channel，将它赋值给conn的的消息接收channel
}

type ServerState int32

const (
	SERVER_STATUS_LAUNCHING = ServerState(iota)
	SERVER_STATUS_RUNNING
	SERVER_STATUS_STOPPING
)

type Option func(*ServerBase)

func WithGenSession(genSession genSession) Option {
	return func(o *ServerBase) {
		o.genSession = genSession
	}
}

func WithReadMsgFunc(readMsgFunc readMsgFunc) Option {
	return func(o *ServerBase) {
		o.readMsgFunc = readMsgFunc
	}
}

func WithConTTLFunc(duration time.Duration) Option {
	return func(o *ServerBase) {
		o.conTTL = duration
	}
}

type ServerBase struct {
	svrType   int32
	svrID     int32
	closeChan chan struct{}
	state     *atomic.Int32
	WhiteList
	Sessions   map[int32]*BaseSession
	MsgChannel chan INetMsg
	conTTL     time.Duration // 连接生效时间（每次接收或发送消息时是顺时延长）
	genSession
	readMsgFunc
}

func (s *ServerBase) SetState(state ServerState) {
	s.state.Store(int32(state))
}

func newServerBase(ops ...Option) ServerBase {
	serverBase := &ServerBase{
		closeChan:  make(chan struct{}, 1),
		state:      atomic.NewInt32(int32(SERVER_STATUS_LAUNCHING)),
		MsgChannel: make(chan INetMsg, 512),
		conTTL:     1 * time.Minute,
	}

	for _, op := range ops {
		op(serverBase)
	}

	if serverBase.readMsgFunc == nil {
		serverBase.readMsgFunc = defaultReadMsg
	}

	if serverBase.genSession == nil {
		serverBase.genSession = newBaseSession
	}

	return *serverBase
}

func (s *ServerBase) AfterClose(callback func()) {
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

				if callback != nil {
					callback()
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

func (s *ServerBase) GetReceiveMsgChan() chan INetMsg {
	return s.MsgChannel
}
