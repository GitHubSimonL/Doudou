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
