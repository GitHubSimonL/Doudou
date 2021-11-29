package network

import (
	"Doudou/lib/logger"
	"net"
	"runtime/debug"
)

type ServerReceiveConnect func(conn net.Conn)

// Server抽象接口
type IServer interface {
	Close()                                                                                // 关闭
	AfterClose(trigger func())                                                             // 关闭
	GetType() int32                                                                        // 类型
	GetID() int32                                                                          // ID
	LoadWhiteList(filename string) bool                                                    // 加载白名单
	AccessCheck(ip string) bool                                                            // 是否放行
	StartListen(port string, rDeadLine, wDeadLine int32, receiveFunc ServerReceiveConnect) // 开始监听
}

type ServerBase struct {
	svrType   int32
	svrID     int32
	closeChan chan struct{}
	WhiteList
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

func (s *ServerBase) GetID() int32 {
	return s.svrID
}
