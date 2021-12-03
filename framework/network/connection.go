package network

import (
	"Doudou/framework/itr"
	. "Doudou/framework/network/default"
	"Doudou/lib/logger"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var pool sync.Pool

type Connection struct {
	Server      itr.IServer            // 当前Conn属于哪个Server
	net.Conn                           // 当前连接的socket TCP套接字
	connID      uint32                 // 当前连接的ID 也可以称作为SessionID，ID全局唯一
	ApiMgr      itr.IApiMgr            // 消息管理MsgID和对应处理方法的消息管理模块
	msgChan     chan []byte            // 无缓冲管道
	msgBuffChan chan []byte            // 有缓冲管道
	isClosed    bool                   // 当前连接的关闭状态
	ctx         context.Context        // 用于控制消息发送与接收协程间同步链接停止
	cancel      context.CancelFunc     // 用于控制消息发送与接收协程间同步链接停止
	metaLock    sync.Mutex             // 保护当前meta的锁
	meta        map[string]interface{} // 链接属性
	sync.RWMutex
	startOnce  sync.Once
	endOnce    sync.Once
	packet     itr.IPacket
	closSignal chan struct{}
}

var _ itr.IConnection = (*Connection)(nil)

func (c *Connection) Start() {
	c.startOnce.Do(func() {
		c.ctx, c.cancel = context.WithCancel(context.Background())

		go c.ReaderTaskStart()
		go c.WriterTaskStart()

		if c.Server != nil {
			c.Server.CallConnStartHookFunc(c)
		}
	})
}

func (c *Connection) Stop() {
	c.endOnce.Do(func() {
		c.Lock()
		defer c.Unlock()

		if c.isClosed {
			return
		}

		defer func() {
			c.cancel()
			c.Close()

			close(c.msgBuffChan)
			close(c.msgChan)
			c.isClosed = true
			c.closSignal <- struct{}{}
		}()

		if c.Server != nil {
			c.Server.CallConnEndHookFunc(c)
			c.Server.GetConnMgr().Remove(c)
		}
	})
}

func (c *Connection) GetContext() context.Context {
	return c.ctx
}

func (c *Connection) GetConnID() uint32 {
	return c.connID
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed {
		return errors.New("Connection is closed!")
	}

	msg, err := c.packet.Pack(NewMessage(msgID, data))
	if err != nil {
		logger.LogErrf("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed {
		return errors.New("Connection is closed!")
	}

	msg, err := c.packet.Pack(NewMessage(msgID, data))
	if err != nil {
		logger.LogErrf("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	c.msgBuffChan <- msg
	return nil
}

func (c *Connection) SetMeta(key string, value interface{}) {
	c.metaLock.Lock()
	defer c.metaLock.Unlock()

	if c.meta == nil {
		c.meta = make(map[string]interface{})
	}

	c.meta[key] = value
}

func (c *Connection) GetMeta(key string) (interface{}, error) {
	c.metaLock.Lock()
	defer c.metaLock.Unlock()
	value, ok := c.meta[key]
	if !ok {
		return nil, fmt.Errorf("Meta key %v not found!", key)
	}

	return value, nil
}

func (c *Connection) RemoveMeta(key string) {
	c.metaLock.Lock()
	defer c.metaLock.Unlock()

	delete(c.meta, key)
}

func (c *Connection) IsClosed() bool {
	c.RLock()
	defer c.RUnlock()

	return c.isClosed
}

func (c *Connection) CloseSignal() chan struct{} {
	return c.closSignal
}

// 生成一个链接对象 (当为client链接时，server对象为空)
func NewConnection(server itr.IServer, conn net.Conn, connID uint32, msgBufferLen int, apiMgr itr.IApiMgr, packet itr.IPacket) *Connection {
	if conn == nil {
		return nil
	}

	return &Connection{
		Server:      server,
		Conn:        conn,
		connID:      connID,
		ApiMgr:      apiMgr,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, msgBufferLen),
		RWMutex:     sync.RWMutex{},
		meta:        make(map[string]interface{}),
		metaLock:    sync.Mutex{},
		isClosed:    false,
		packet:      packet,
		closSignal:  make(chan struct{}, 1),
	}

}

// 写任务开启
func (c *Connection) WriterTaskStart() {
	logger.LogDebugf("Conn:%v Writer Goroutine is running!", c.GetConnID())
	defer logger.LogDebugf("Conn:%v Remote:%v Reader exit!", c.GetConnID(), c.RemoteAddr().String())
	defer c.Stop()

	for {
		select {
		case data, ok := <-c.msgChan:
			if !ok {
				return
			}

			if _, err := c.Write(data); err != nil {
				logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
			}

		case data, ok := <-c.msgBuffChan:
			if !ok {
				logger.LogDebugf("msgBuffChan is Closed")
				return
			}

			if _, err := c.Write(data); err != nil {
				logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

// 读任务开启
func (c *Connection) ReaderTaskStart() {
	logger.LogDebugf("Conn:%v Reader Goroutine is running!", c.GetConnID())
	defer logger.LogDebugf("Conn:%v Remote:%v Reader exit!", c.GetConnID(), c.RemoteAddr().String())
	defer c.Stop()

	br := bufio.NewReader(c)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:

			// 获取包头
			header := make([]byte, c.packet.GetHeadLen())
			if _, err := io.ReadFull(br, header); err != nil {
				if err != io.EOF {
					logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
				}
				return
			}

			// 解包
			head, err := c.packet.UnpackHead(header)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}

			data := make([]byte, head.GetDataLen())
			if head.GetDataLen() > 0 {
				if _, err := io.ReadFull(br, data); err != nil {
					if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
						logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
					}
					return
				}
			}

			msg := NewMessage(head.GetMsgID(), data)
			c.SetDeadline(time.Now().Add(DefaultConnectionTTL))

			req := &Request{
				conn: c,
				msg:  msg,
			}

			if c.ApiMgr.GetTaskQueueAmount() > 1 { // 多任务处理
				c.ApiMgr.AddMgsToTaskPool(req)
				continue
			}

			// 单协程处理
			c.ApiMgr.DoMsgHandler(req)
		}
	}
}
