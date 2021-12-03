package itr

type IHead interface {
	SetMsgID(msgID uint32)
	SetDataLen(dataLen int)
	GetMsgID() uint32
	GetDataLen() int
	GetHeadLen() int
}

type Head struct {
	MsgID   uint32
	DataLen int
}

var _ IHead = (*Head)(nil)

func (h *Head) SetMsgID(msgID uint32) {
	h.MsgID = msgID
}

func (h *Head) SetDataLen(dataLen int) {
	h.DataLen = dataLen
}

func (h *Head) GetMsgID() uint32 {
	return h.MsgID
}

func (h *Head) GetDataLen() int {
	return h.DataLen
}

func (h *Head) GetHeadLen() int {
	return 8
}
