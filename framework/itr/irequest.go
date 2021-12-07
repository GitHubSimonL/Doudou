package itr

type IRequest interface {
	GetConnection() IConnection // 获取请求连接信息
	GetData() interface{}       // 获取请求消息的数据
	GetMsgID() uint32           // 获取请求的消息ID
}
