package network

import (
	"net"
)

type Session struct {
	SessionID int64
	Conn      net.Conn
}
