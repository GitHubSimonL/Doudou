package _default

import (
	"Doudou/framework/itr"
	"bytes"
	"encoding/binary"
)

type NetPack struct {
	itr.IHead
}

func (n *NetPack) UnpackHead(binaryData []byte) (itr.IHead, error) {
	head := &itr.Head{
		MsgID:   0,
		DataLen: 0,
	}

	dataBuff := bytes.NewReader(binaryData)

	if err := binary.Read(dataBuff, binary.LittleEndian, head); err != nil {
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

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (n *NetPack) GetHeadLen() int {
	return n.GetHeadLen()
}
