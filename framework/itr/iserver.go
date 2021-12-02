package itr

type IServer interface {
	Start()                                      // 启动服务器
	Stop()                                       // 停止
	SetHandler(msgID int32, handle IHandle)      // 根据MsgID设置handle方法
	GetConnMgr() IConnMgr                        // 获取server所有链接管理器
	SetConnStartHookFunc(func(conn IConnection)) // 链接创建时的hood方法
	CallConnStartHookFunc(conn IConnection)      // 调用链接创建hood方法
	SetConnEndHookFunc(func(conn IConnection))   // 链接断开时的hood方法
	CallConnEndHookFunc(conn IConnection)        // 调用链接断开hood方法
	Packet() Packet
}