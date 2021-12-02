package framework

import (
	"Doudou/lib/logger"
	"bufio"
	"time"
)

var curSessionID uint32 = 1

const ConnTimeOut = 1 * time.Minute // 链接超时时间（每次接收消息后，会将超时设置成2分钟过后）

type ISession interface {
	sendMsg(data []byte)               // 发送消息
	SetSvrReceiveMsgChan(chan INetMsg) // 设置住
	SendMsg(msg INetMsg)               // 发送消息
	Start()                            // 开始工作
	IsClosed() bool                    // 是否关闭
	Close()                            // 关闭
	SetReadMsgFunc(readMsgFunc)        // 设置读消息方法
	GetSessionID() uint32
}

type BaseSession struct {
	sessionID      uint32
	userID         int64
	conn           ICon
	receiveMsgChan chan INetMsg // 这个直接和server共用一个channel
	sendMsgChan    chan INetMsg
	isClosed       bool
	ttl            time.Duration
	readMsgFunc
	sender *MsgSender
}

func (b *BaseSession) BindUserID(userID int64) {
	b.userID = userID
}

func (b *BaseSession) GetSessionID() uint32 {
	return b.sessionID
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

func (b *BaseSession) SetSvrReceiveMsgChan(receiveMsgChan chan INetMsg) {
	if receiveMsgChan == nil {
		return
	}

	b.receiveMsgChan = receiveMsgChan
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

	b.conn.Write(data)

	isSuccess = true
	return
}

func (b *BaseSession) Start() {
	b.sender.start()

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

			b.conn.SetDeadline(time.Now().Add(b.ttl))

			netMsg := b.readMsgFunc(rd)
			if netMsg == nil {
				return
			}

			if b.userID > 0 && netMsg.GetUserID() == 0 {
				netMsg.SetUserID(b.userID)
			}

			if b.userID == 0 && netMsg.GetUserID() > 0 {
				b.BindUserID(b.userID)
			}

			netMsg.SetSessionID(b.GetSessionID())

			b.receiveMsgChan <- netMsg
		}
	}()
}

func newBaseSession(conn ICon) ISession {
	if conn == nil {
		return nil
	}

	defer func() {
		curSessionID++
	}()

	session := &BaseSession{
		sessionID: curSessionID,
		conn:      conn,
		ttl:       ConnTimeOut,
		sender:    newMsgSender(conn, 0),
	}

	return session
}