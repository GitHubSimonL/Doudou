package itr

type IConnMgr interface {
	Add(conn IConnection)                       // 添加链接
	Remove(conn IConnection)                    // 删除连接(这里是外部调用，需要调用方保证conn被正确close)
	Get(connID uint32) (IConnection, error)     // 利用ConnID获取链接
	GetByAddr(addr string) (IConnection, error) // 根据地址获取 （udp）
	Len() int                                   // 获取当前连接
	ClearConn()                                 // 删除并停止所有链接
}
