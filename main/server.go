package main

import (
	"earnth/eiface"
	"earnth/enet"
	"fmt"
)

type PingRouter struct {
	enet.BaseRouter
}

// Handle test
func (pr *PingRouter) Handle(request eiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("PingRouter recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回写数据
	err := request.GetConnection().SendMsg(0, []byte("123456"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloEarnthRouter struct {
	enet.BaseRouter
}

func (hr *HelloEarnthRouter) Handle(request eiface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("HelloEarnthRouter recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建连接时执行
func DoConnectionBegin(conn eiface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn eiface.IConnection) {
	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	s := enet.NewServer()
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloEarnthRouter{})
	s.Serve()
}
