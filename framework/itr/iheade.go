package itr

type IHead interface {
	SetMsgID(msgID uint32)
	SetDataLen(dataLen int32)
	GetMsgID() uint32
	GetDataLen() int32
	GetHeadLen() int32
	NewHead() IHead
}

type Head struct {
	MsgID   uint32
	DataLen int32
}

func (h *Head) NewHead() IHead {
	return &Head{}
}

var _ IHead = (*Head)(nil)

func (h *Head) SetMsgID(msgID uint32) {
	h.MsgID = msgID
}

func (h *Head) SetDataLen(dataLen int32) {
	h.DataLen = dataLen
}

func (h *Head) GetMsgID() uint32 {
	return h.MsgID
}

func (h *Head) GetDataLen() int32 {
	return int32(h.DataLen)
}

func (h *Head) GetHeadLen() int32 {
	return int32(8)
}
