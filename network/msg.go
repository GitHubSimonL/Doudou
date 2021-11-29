package network

import (
	"errors"
	"io"
	"net"
	"time"
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

const (
	NetMessageHeaderSize = 13
	HeaderFlagCrypto     = 0x02
	HeaderFlagCompress   = 0x08
	HeaderFlagCrc        = 0x20

	MaxHandShakeMsgLen    = 1024 * 10
	MaxExternalMessageLen = 1024 * 1024 * 2  // 对client的最大的消息长度
	MaxInternalMessageLen = 1024 * 1024 * 64 // server内部最大的消息长度
)

type MsgPackHeader struct {
	Length uint32
	Flag   byte
	MsgID  uint32
	CRC    uint32
}

const ConnTimeOut = 2 * time.Minute // 链接超时时间（每次接收消息后，会将超时设置成2分钟过后）

type INetMsg interface {
	Encode() (bData []byte)
	GetSessionID() uint32   // 获取消息的会话ID
	SetSessionID(id uint32) // 消息获得后，设置来自那个session，方便处理逻辑后回包
}

type NetMsg struct {
	Time     int64
	MsgId    uint32
	Payload  interface{}
	ClientId int64 // 客户端id，如果连接多个同类服务器，记录消息来自哪个连接
}

// 读消息
func defaultReadMsg(conn net.Conn, rd io.Reader) INetMsg {
	defer func() {
		conn.SetReadDeadline(time.Now().Add(ConnTimeOut))
	}()
	return nil
}
