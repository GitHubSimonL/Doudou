package framework

import "net"

type IClient interface {
	Send(msg INetMsg)
}

type BaseClient struct {
	conn   ICon
	sender *MsgSender
}

func (b *BaseClient) Send(msg INetMsg) error {
	if b.sender == nil || msg == nil {
		return nil
	}

	return b.sender.Send(msg)
}

type TcpClient struct {
	BaseClient
}

func NewTcpClient(addr string) (*TcpClient, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	dCon := newConnection(conn)

	tc := &TcpClient{
		BaseClient: BaseClient{
			conn:   dCon,
			sender: newMsgSender(dCon, 0),
		},
	}

	tc.sender.start()
	return tc, nil
}
