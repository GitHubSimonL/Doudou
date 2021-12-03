package itr

type IPacket interface {
	Unpack(head IHead, binaryData []byte) (IMessage, error)
	UnpackHead(binaryData []byte) (IHead, error)
	Pack(msg IMessage) ([]byte, error)
	GetHeadLen() int32
}
