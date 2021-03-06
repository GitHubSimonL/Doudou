package itr

import (
	"context"
	"net"
)

// 定义连接接口
type IConnection interface {
	Start()                                      // 启动连接，让当前连接开始工作
	Stop()                                       // 停止连接，结束当前连接状态
	GetContext() context.Context                 // 返回ctx，用于用户自定义的go程获取连接退出状态
	GetConnID() uint32                           // 获取当前连接ID
	SendMsg(msgID uint32, data []byte) error     // 直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	SendBuffMsg(msgID uint32, data []byte) error // 直接将Message数据发送给远程的TCP客户端(有缓冲)
	SetMeta(key string, value interface{})       // 设置链接属性
	GetMeta(key string) (interface{}, error)     // 获取链接属性
	RemoveMeta(key string)                       // 移除链接属性
	IsClosed() bool                              // 链接是否已关闭
	CloseSignal() chan struct{}
	net.Conn
}
