package itr

type Option func(server IServer)

type IServer interface {
	Start()                                      // 启动服务器
	Stop()                                       // 停止
	SetHandler(msgID int32, handle IHandle)      // 根据MsgID设置handle方法
	GetConnMgr() IConnMgr                        // 获取server所有链接管理器
	SetConnStartHookFunc(func(conn IConnection)) // 链接创建时的hook方法
	CallConnStartHookFunc(conn IConnection)      // 调用链接创建hook方法
	SetConnEndHookFunc(func(conn IConnection))   // 链接断开时的hook方法
	CallConnEndHookFunc(conn IConnection)        // 调用链接断开hook方法
	Packet() IPacket                             // 数据打包与解包对象
	SetType(svrType int32)                       // 设置类型
	SetID(svrID int32)                           // 设置ID
	SetIP(ip string)                             // 设置IP
	SetPort(port int)                            // 设置端口
}

func WithSvrType(svrType int32) Option {
	return func(server IServer) {
		server.SetType(svrType)
	}
}

func WithSvrID(svrID int32) Option {
	return func(server IServer) {
		server.SetID(svrID)
	}
}

func WithIP(ip string) Option {
	return func(server IServer) {
		server.SetIP(ip)
	}
}

func WithPort(port int) Option {
	return func(server IServer) {
		server.SetPort(port)
	}
}
