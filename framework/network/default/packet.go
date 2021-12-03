package _default

import (
	"Doudou/framework/itr"
	"bytes"
	"encoding/binary"
)

type NetPack struct {
	itr.IHead
}

var _ itr.IPacket = (*NetPack)(nil)

func (n *NetPack) UnpackHead(binaryData []byte) (itr.IHead, error) {
	head := &itr.Head{
		MsgID:   0,
		DataLen: 0,
	}

	dataBuff := bytes.NewReader(binaryData)

	if err := binary.Read(dataBuff, binary.LittleEndian, &head.MsgID); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &head.DataLen); err != nil {
		return nil, err
	}
	return head, nil
}

func NewNetPacket() itr.IPacket {
	return &NetPack{
		IHead: &itr.Head{},
	}
}

func (n *NetPack) Unpack(head itr.IHead, binaryData []byte) (itr.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	msg.SetMsgID(head.GetMsgID())

	return msg, nil
}

func (n *NetPack) Pack(msg itr.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	head := n.NewHead()
	head.SetDataLen(msg.GetDataLen())
	head.SetMsgID(msg.GetMsgID())

	if err := binary.Write(dataBuff, binary.LittleEndian, head.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, head.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (n *NetPack) GetHeadLen() int32 {
	return n.IHead.GetHeadLen()
}
