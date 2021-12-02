package framework

import (
	"errors"
	"io"
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

type INetMsg interface {
	Encode() (bData []byte)
	GetSessionID() uint32   // 获取消息的会话ID
	SetSessionID(id uint32) // 消息获得后，设置来自那个session，方便处理逻辑后回包
	SetUserID(userID int64) // 设置玩家ID
	GetUserID() int64
}

type DefaultMsg struct {
	Data      []byte
	sessionID uint32
	userID    int64
}

func (n *DefaultMsg) Encode() (bData []byte) {
	return n.Data
}

func (n *DefaultMsg) GetSessionID() uint32 {
	return n.sessionID
}

func (n *DefaultMsg) SetSessionID(id uint32) {
	n.sessionID = id
}

func (n *DefaultMsg) SetUserID(userID int64) {
	n.userID = userID
}

func (n *DefaultMsg) GetUserID() int64 {
	return n.userID
}

// 读消息
func defaultReadMsg(rd io.Reader) INetMsg {
	data := make([]byte, 20)
	_, err := io.ReadFull(rd, data)
	if err != nil {
		return nil
	}

	return &DefaultMsg{Data: data}
}
