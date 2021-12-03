package _default

import "Doudou/framework/itr"

type Head struct {
	MsgID   uint32
	DataLen int
}

type Message struct {
	MsgID uint32
	Data  []byte
}

var _ itr.IMessage = (*Message)(nil)

func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		MsgID: msgID,
		Data:  data,
	}
}

func (m *Message) GetDataLen() int {
	return len(m.Data)
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
