package network

import (
	"Doudou/lib/logger"
	"encoding/binary"
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

const ConnTimeOut = 2 * time.Minute

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

// 读消息
func readMsg(conn net.Conn, rd io.Reader, maxMsgBodyLen uint32) ([]byte, *MsgPackHeader) {
	defer func() {
		conn.SetReadDeadline(time.Now().Add(ConnTimeOut))
	}()

	headerBuff := make([]byte, NetMessageHeaderSize)
	_, err := io.ReadFull(rd, headerBuff)

	errF := func(err error) {
		if errors.Is(err, io.EOF) {
			logger.Logf("Recv EOF: %v.", conn.RemoteAddr())
		} else {
			logger.LogWarnf("Recv err: %v %v.", conn.RemoteAddr(), err)
		}
	}

	if err != nil {
		errF(err)
		return nil, nil
	}

	header := &MsgPackHeader{
		Length: binary.BigEndian.Uint32(headerBuff),
		Flag:   headerBuff[4],
		MsgID:  binary.BigEndian.Uint32(headerBuff[5:]),
		CRC:    binary.BigEndian.Uint32(headerBuff[9:]),
	}

	if header.Length < NetMessageHeaderSize || header.Length > maxMsgBodyLen {
		logger.LogWarnf("message %d len %d out of range", header.MsgID, header.Length)
		return nil, nil
	}

	data := make([]byte, header.Length-NetMessageHeaderSize)
	_, err = io.ReadFull(rd, data)
	if err != nil {
		logger.LogWarnf("Recv msg %d from %v, error: %v.", header.MsgID, conn.RemoteAddr(), err)
		return nil, nil
	}

	return data, header
}
