package network

import (
	"Doudou/lib/logger"
	"bufio"
	"net"
)

var curSessionID uint32 = 1

type ISession interface {
	receiveMsg()                       // 接收消息
	sendMsg(data []byte)               // 发送消息
	SetSvrReceiveMsgChan(chan INetMsg) // 设置住
	SendMsg(msg INetMsg)               // 发送消息
	Start()                            // 开始工作
	IsClosed() bool                    // 是否关闭
	Close()                            // 关闭
	SetReadMsgFunc(readMsgFunc)        // 设置读消息方法
}

type BaseSession struct {
	SessionID uint32
	conn      net.Conn
	msgChan   chan INetMsg
	isClosed  bool
	readMsgFunc
}

func (b *BaseSession) receiveMsg() {
	panic("implement me")
}

func (b *BaseSession) sendMsg(data []byte) {
	panic("implement me")
}

func (b *BaseSession) SetReadMsgFunc(fnc readMsgFunc) {
	if fnc == nil {
		return
	}

	b.readMsgFunc = fnc
}

func (b *BaseSession) Close() {
	if b.conn == nil {
		return
	}

	defer func() {
		b.isClosed = true
	}()

	b.conn.Close()
}

func (b *BaseSession) IsClosed() bool {
	if b.conn == nil {
		return true
	}

	return b.isClosed
}

func (b *BaseSession) SetSvrReceiveMsgChan(msgChan chan INetMsg) {
	if msgChan == nil {
		return
	}

	b.msgChan = msgChan
}

func (b *BaseSession) SendMsg(msg INetMsg) {
	if b.IsClosed() {
		return
	}

	var isSuccess bool
	defer func() {
		if !isSuccess {
			b.Close()
			return
		}
	}()

	data := msg.Encode()
	if len(data) <= 0 {
		return
	}

	isSuccess = true
	return
}

func (b *BaseSession) Start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogErrf("server session panic: %v.", r)
			}

			b.Close()
		}()

		if b.readMsgFunc == nil {
			return
		}

		var rd = bufio.NewReader(b.conn)
		for {
			if b.IsClosed() {
				return
			}

			netMsg := b.readMsgFunc(b.conn, rd)
			if netMsg == nil {
				return
			}

			b.msgChan <- netMsg
		}
	}()
}

func newBaseSession(conn net.Conn) ISession {
	if conn == nil {
		return nil
	}

	defer func() {
		curSessionID++
	}()

	session := &BaseSession{
		SessionID: curSessionID,
		conn:      conn,
	}

	return session
}
