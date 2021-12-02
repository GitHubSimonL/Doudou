package nerwork

import (
	"Doudou/framework/itr"
)

type Request struct {
	conn itr.IConnection
	msg  itr.IMessage
}

func (r *Request) GetConnection() itr.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
