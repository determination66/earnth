package enet

import (
	"earnth/demo1/eiface"
	"earnth/demo1/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前conn属于的Server
	TcpServer eiface.IServer

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
	//有缓冲的管道，用于读写两个goroutine之间的通信
	msgBuffChan chan []byte

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewConnection(server eiface.IServer, conn *net.TCPConn, connID uint32, msgHandler eiface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClose:     false,
		MsgHandler:  msgHandler,
		ExitChan:    make(chan bool, 1),
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]interface{}), //对链接属性map初始化
	}
	//添加到链接管理当中
	c.TcpServer.GetConnManager().Add(c)

	return c
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
		//go c.MsgHandler.DoMsgHandler(req)
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
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
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
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
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)

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

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()
	c.ExitChan <- true

	//删除当前的链接
	c.TcpServer.GetConnManager().Remove(c)

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
		return errors.New("connection closed when send msg,SendMsg")
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

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClose {
		return errors.New("connection closed when send msg,SendBuffMsg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
