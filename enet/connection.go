package enet

import (
	"earnth/eiface"
	"net"
)

type Connection struct {
	// 当前tcp套接字
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//当前连接状态
	isClose bool

	//当前连接所绑定的业务方法API
	handleAPI eiface.HandleFunc

	//告知当前连接已经退出
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI eiface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClose:   false,
		handleAPI: callbackAPI,
		ExitChan:  make(chan bool, 1),
	}
}
