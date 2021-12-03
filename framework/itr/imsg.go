package itr

type IMessage interface {
	GetDataLen() int32 //	获取消息数据段长度
	GetMsgID() uint32  //	获取消息ID
	GetData() []byte   //	获取消息内容
	SetMsgID(uint32)   //	设计消息ID
	SetData([]byte)    //	设计消息内容
}
