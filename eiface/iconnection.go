package eiface

import "net"

// IConnection 定义连接模块的接口
type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止连接
	Stop()
	// GetTcpConnection 获取当前连接的socket conn
	GetTcpConnection() *net.TCPConn
	//GetConnId 获取远程的tcp的ID
	GetConnId() uint32
	// RemoteAddr 获取远程客户端TCP状态
	RemoteAddr() net.Addr
	// Send 发送数据，发送给客户端
	Send(data []byte) error
	//SendMsg 直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error
	// SendBuffMsg 带缓冲的发送
	SendBuffMsg(msgId uint32, data []byte) error
}

// HandleFunc 业务处理
type HandleFunc func(*net.TCPConn, []byte, int) error
