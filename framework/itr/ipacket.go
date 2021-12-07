package itr

type IPacket interface {
	UnpackHead(binaryData []byte) (IHead, error)
	Pack(msg IMessage) ([]byte, error)
	GetHeadLen() int32
	Unpack2IRequest(conn IConnection, msgID uint32, binaryData []byte) (IRequest, error)
}
