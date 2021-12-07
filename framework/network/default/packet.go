package _default

import (
	"Doudou/framework/itr"
	"bytes"
	"encoding/binary"
	"errors"
)

type NetPack struct {
	itr.IHead
}

var _ itr.IPacket = (*NetPack)(nil)

func NewNetPacket() itr.IPacket {
	return &NetPack{
		IHead: &itr.Head{},
	}
}

func (n *NetPack) UnpackHead(binaryData []byte) (itr.IHead, error) {
	head := &itr.Head{
		MsgID:   0,
		DataLen: 0,
	}

	dataBuff := bytes.NewReader(binaryData)

	if err := binary.Read(dataBuff, binary.LittleEndian, &head.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &head.MsgID); err != nil {
		return nil, err
	}

	return head, nil
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

func (n *NetPack) GetHeadLen() int32 {
	return n.IHead.GetHeadLen()
}

func (n *NetPack) Unpack2IRequest(conn itr.IConnection, msgID uint32, binaryData []byte) (itr.IRequest, error) {
	if conn == nil {
		return nil, errors.New("param is nil")
	}

	return &Request{
		conn: conn,
		msg: &Message{
			MsgID: msgID,
			Data:  binaryData,
		},
	}, nil
}
