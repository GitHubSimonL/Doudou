package nerwork

import (
	"Doudou/framework/itr"
	"Doudou/lib/logger"
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"sync"
)

var pool sync.Pool

type Connection struct {
	Server      itr.IServer            // 当前Conn属于哪个Server
	conn        net.Conn               // 当前连接的socket TCP套接字
	connID      uint32                 // 当前连接的ID 也可以称作为SessionID，ID全局唯一
	MsgHandler  itr.IHandleMgr         // 消息管理MsgID和对应处理方法的消息管理模块
	msgChan     chan []byte            // 无缓冲管道
	msgBuffChan chan []byte            // 有缓冲管道
	isClosed    bool                   // 当前连接的关闭状态
	ctx         context.Context        // 用于控制消息发送与接收协程间同步链接停止
	cancel      context.CancelFunc     // 用于控制消息发送与接收协程间同步链接停止
	metaLock    sync.Mutex             // 保护当前meta的锁
	meta        map[string]interface{} // 链接属性
	sync.RWMutex
}

var _ itr.IConnection = (*Connection)(nil)

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.ReaderTaskStart()
	go c.WriterTaskStart()

	c.Server.CallConnStartHookFunc(c)
}

func (c *Connection) Stop() {
	c.Lock()
	defer c.Unlock()

	if c.isClosed {
		return
	}

	defer func() {
		c.cancel()
		c.GetConn().Close()

		close(c.msgBuffChan)
		close(c.msgChan)
	}()

	c.Server.CallConnEndHookFunc(c)
	c.Server.GetConnMgr().Remove(c)
}

func (c *Connection) GetContext() context.Context {
	return c.ctx
}

func (c *Connection) GetConn() net.Conn {
	return c.GetConn()
}

func (c *Connection) GetConnID() uint32 {
	return c.connID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed {
		return errors.New("Connection is closed!")
	}

	sp := c.Server.Packet()
	msg, err := sp.Pack(NewMessage(msgID, data))
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

	sp := c.Server.Packet()
	msg, err := sp.Pack(NewMessage(msgID, data))
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

func (c *Connection) RemoveProperty(key string) {
	panic("implement me")
}

func (c *Connection) IsClosed() bool {
	panic("implement me")
}

// 生成一个链接对象
func NewConnection(server itr.IServer, conn net.Conn, connID uint32, msgHandleMgr itr.IHandleMgr, msgBufferLen int) *Connection {
	if server == nil || conn == nil || msgHandleMgr == nil {
		return nil
	}

	return &Connection{
		Server:      server,
		conn:        conn,
		connID:      connID,
		MsgHandler:  msgHandleMgr,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, msgBufferLen),
		RWMutex:     sync.RWMutex{},
		meta:        make(map[string]interface{}),
		metaLock:    sync.Mutex{},
		isClosed:    false,
	}

}

// 写任务开启
func (c *Connection) WriterTaskStart() {
	logger.LogDebugf("Conn:%v Writer Goroutine is running!", c.GetConnID())
	defer logger.LogDebugf("Conn:%v Remote:%v Reader exit!", c.GetConnID(), c.RemoteAddr().String())

	for {
		select {
		case data, ok := <-c.msgChan:
			if !ok {
				return
			}

			if _, err := c.GetConn().Write(data); err != nil {
				logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
			}

		case data, ok := <-c.msgBuffChan:
			if !ok {
				logger.LogDebugf("msgBuffChan is Closed")
				return
			}

			if _, err := c.GetConn().Write(data); err != nil {
				logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
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

	br := bufio.NewReader(c.GetConn())

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			header := pool.Get().([]byte)
			if header == nil {
				header = make([]byte, c.Server.Packet().GetHeadLen())
			}

			// 获取包头
			if _, err := io.ReadFull(br, header); err != nil {
				if err != io.EOF {
					logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
				}
				return
			}

			// 解包
			msg, err := c.Server.Packet().Unpack(header)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}

			pool.Put(header)

			data := pool.Get().([]byte)
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.GetConn(), data); err != nil {
					if err != io.EOF {
						logger.LogErrf("Conn:%v Remote:%v catch err.%v", c.GetConnID(), c.RemoteAddr().String(), err.Error())
					}
					return
				}
			}

			msg.SetData(data)
			req := &Request{
				conn: c,
				msg:  msg,
			}

			if c.MsgHandler.GetTaskQueueAmount() > 1 { // 多任务处理
				c.MsgHandler.AddMgsToTaskPool(req)
				continue
			}

			// 单协程处理
			c.MsgHandler.DoMsgHandler(req)
		}
	}
}
