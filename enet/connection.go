package enet

import (
	"earnth/eiface"
	"fmt"
	"net"
	"time"
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

func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println("ConnId=", c.ConnID, "reader exit ,remote Addr is:", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端数据到buf
		buf := make([]byte, 512)
		var err error
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			time.Sleep(time.Second)
			break
		}
		//调用当前连接绑定的HandlAPI
		err = c.handleAPI(c.Conn, buf, cnt)
		if err != nil {
			fmt.Println("ConnID", c.ConnID, " handle err:", err)
			break
		}

	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID", c.ConnID)
	//启动从当前读数据业务
	//todo 启动写数据业务
	go c.StartReader()

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID =", c.ConnID)
	if c.isClose == true {
		return
	}
	c.isClose = true

	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}
