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
