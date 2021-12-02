package _default

import "Doudou/framework/itr"

type Message struct {
	Len   uint32
	MsgID uint32
	Data  []byte
}

var _ itr.IMessage = (*Message)(nil)

func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		Len:   uint32(len(data)),
		MsgID: msgID,
		Data:  data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.Len
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(msgID uint32) {
	m.MsgID = msgID
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(msgLen uint32) {
	m.Len = msgLen
}
