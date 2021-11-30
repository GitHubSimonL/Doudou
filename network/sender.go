package network

import (
	"Doudou/lib/logger"
	"fmt"
)

// 消息发送器
type MsgSender struct {
	conn       ICon
	exitSignal chan struct{}
	msgList    chan INetMsg
}

// 异步发送数据
func (ms *MsgSender) Send(msg INetMsg) error {
	if msg == nil || ms.conn == nil || ms.conn.IsClosed() {
		return nil
	}

	select {
	case ms.msgList <- msg:
		return nil
	default:
		err := fmt.Errorf("sender overflow, pending %d, remote: %v", len(ms.msgList), ms.conn.RemoteAddr())
		logger.LogWarnf("%v", err)
		return err
	}
}

// 同步发送数据
func (ms *MsgSender) SyncSend(msg INetMsg) (err error) {
	return ms.rawSend(msg)
}

func (ms *MsgSender) rawSend(msg INetMsg) error {
	if msg == nil || ms.conn == nil || ms.conn.IsClosed() {
		return nil
	}

	_, err := ms.conn.Write(msg.Encode())
	return err
}

// 发送线程启动
func (ms *MsgSender) start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogErrf("Sender %v error: %v.", ms, r)
			}
		}()

		for {
			select {
			case msg := <-ms.msgList:
				if err := ms.rawSend(msg); err != nil {
					logger.LogErrf("send message %v: %v", msg, err)
				}

			case <-ms.exitSignal:
				logger.Logf("Sender close: %p %v.", ms, ms)
				return
			}
		}
	}()
}
