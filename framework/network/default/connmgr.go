package _default

import (
	"Doudou/framework/itr"
	"fmt"
	"sync"
)

type ConnMgr struct {
	connectionsMap map[uint32]itr.IConnection
	sync.RWMutex
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connectionsMap: make(map[uint32]itr.IConnection),
		RWMutex:        sync.RWMutex{},
	}
}

func (c *ConnMgr) Add(conn itr.IConnection) {
	if conn == nil {
		return
	}

	c.Lock()
	defer c.Unlock()

	if oldConn, ok := c.connectionsMap[conn.GetConnID()]; ok {
		oldConn.Stop()
	}

	c.connectionsMap[conn.GetConnID()] = conn
}

func (c *ConnMgr) Remove(conn itr.IConnection) {
	if conn == nil {
		return
	}

	c.Lock()
	defer c.Unlock()

	delete(c.connectionsMap, conn.GetConnID())
}

func (c *ConnMgr) Get(connID uint32) (itr.IConnection, error) {
	c.RLock()
	defer c.RUnlock()

	conn, ok := c.connectionsMap[connID]
	if !ok {
		return nil, fmt.Errorf("connection %v not found.", connID)
	}

	return conn, nil
}

func (c *ConnMgr) Len() int {
	c.RLock()
	defer c.RUnlock()

	return len(c.connectionsMap)
}

func (c *ConnMgr) ClearConn() {
	c.Lock()
	defer c.Unlock()

	for connID, conn := range c.connectionsMap {
		conn.Stop()

		delete(c.connectionsMap, connID)
	}
}
