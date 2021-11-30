package network

import "net"

type ICon interface {
	IsClosed() bool
	net.Conn
}

// 抽象链接
type DConn struct {
	net.Conn
	isClosed bool
}

var _ net.Conn = (*DConn)(nil)

func (d *DConn) IsClosed() bool {
	return d.isClosed
}

func (d *DConn) Close() error {
	if err := d.Conn.Close(); err != nil {
		return err
	}

	d.isClosed = true
	return nil
}
