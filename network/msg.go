package network

import (
	"errors"
)

var (
	eCrcErr          = errors.New("check crc failed")
	eNeedCryptor     = errors.New("need cryptor")
	eDecryptFailed   = errors.New("decrypt msg failed")
	eMsgLenErr       = errors.New("msg length error")
	eMsgSizeOverflow = errors.New("msg size too big")
	eRecvMsgErr      = errors.New("read data failed")
	eSessionClosed   = errors.New("session is closed")
	eParseMsgFailed  = errors.New("parse msg failed")
)

type INetMsg interface {
	Encode() (bData []byte)
	Decode(bData []byte) NetMsg
}

type NetMsg struct {
	Time     int64
	MsgId    uint32
	Payload  interface{}
	ClientId int64 // 客户端id，如果连接多个同类服务器，记录消息来自哪个连接
}
