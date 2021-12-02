package network

import (
	"Doudou/framework/itr"
	"bytes"
	"encoding/binary"
)

type NetPack struct {
}

func NewNetPacket() itr.IPacket {
	return &NetPack{}
}

func (n *NetPack) Unpack(binaryData []byte) (itr.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgID); err != nil {
		return nil, err
	}

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

func (n *NetPack) GetHeadLen() uint32 {
	return DefaultPackageHeaderLen
}
