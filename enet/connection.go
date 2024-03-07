package enet

import (
	"earnth/eiface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	// 当前tcp套接字
	Conn *net.TCPConn

	//连接ID
	ConnID uint32

	//当前连接状态
	isClose bool

	//告知当前连接已经退出
	ExitChan chan bool

	MsgHandler eiface.IMsgHandler

	//无缓冲管道，用于读写两个goroutine之间的通信
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler eiface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClose:    false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")
	defer fmt.Println("ConnId=", c.ConnID, "reader exit ,remote Addr is:", c.RemoteAddr().String())
	defer c.Stop()

	//创建拆包解包对象
	dp := NewDataPack()
	for {

		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		//拆包
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}

		if msg.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg, ok := msg.(*Message)
			if !ok {
				fmt.Println("msgHead assert failed!")
			}
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(c.GetTcpConnection(), msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Server receive Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}
		//得到当前客户端请求的Request数据
		req := &Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(req)

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			//conn已经关闭
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID", c.ConnID)
	//todo 启动写数据业务
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	for {
		select {
		case <-c.ExitChan:
			//得到退出消息，不再阻塞
			return
		}
	}

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID =", c.ConnID)
	if c.isClose == true {
		return
	}
	c.isClose = true

	c.Conn.Close()
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msgId:", err)
		return err
	}
	//采用管道共享信息
	c.msgChan <- binaryMsg
	return nil
}
