package main

import (
	"earnth/eiface"
	"earnth/enet"
	"fmt"
	"io"
	"net"
)

type PingRouter struct {
	enet.BaseRouter
}

func (pr *PingRouter) PreHandle(req eiface.IRequest) {
	fmt.Println("call func PreHandle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("before ping ...."))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}
func (pr *PingRouter) Handle(req eiface.IRequest) {
	fmt.Println("call func Handle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("ping ping ...."))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}
func (pr *PingRouter) PostHandle(req eiface.IRequest) {
	fmt.Println("call func PostHandle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("after ping ...."))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {
	//s := enet.NewServer()
	//s.AddRouter(&PingRouter{})
	//s.Serve()

	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建服务器gotoutine，负责从客户端goroutine读取粘包的数据，然后进行解析

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}

		//处理客户端请求
		go func(conn net.Conn) {
			//创建封包拆包对象dp
			dp := enet.NewDataPack()
			for {
				//1 先读出流中的head部分
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
				if err != nil {
					fmt.Println("read head error")
					break
				}
				//将headData字节流 拆包到msg中
				msgHead, err := dp.UnPack(headData)
				if err != nil {
					fmt.Println("server unpack err:", err)
					return
				}

				if msgHead.GetDataLen() > 0 {
					//msg 是有data数据的，需要再次读取data数据
					msg := msgHead.(*enet.Message)
					msg.Data = make([]byte, msg.GetDataLen())

					//根据dataLen从io中读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}

					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}
			}
		}(conn)
	}
}
