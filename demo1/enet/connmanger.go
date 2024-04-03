package enet

import (
	"earnth/demo1/eiface"
	"errors"
	"fmt"
	"sync"
)

type ConnManger struct {
	connections map[uint32]eiface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManger {
	return &ConnManger{
		connections: make(map[uint32]eiface.IConnection),
	}
}

func (c *ConnManger) Add(conn eiface.IConnection) {
	//保证安全
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnId()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManger) Remove(conn eiface.IConnection) {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除连接信息
	delete(c.connections, conn.GetConnId())

	fmt.Println("connection Remove ConnID=", conn.GetConnId(), " successfully: conn num = ", c.Len())
}

func (c *ConnManger) Get(connId uint32) (eiface.IConnection, error) {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManger) Len() int {
	return len(c.connections)
}

func (c *ConnManger) CLearConn() {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())
}
