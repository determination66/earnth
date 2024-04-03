package main

import (
	"earnth/demo1/eiface"
	"earnth/demo1/enet"
	"fmt"
)

type PingRouter struct {
	enet.BaseRouter
}

// Handle test
func (pr *PingRouter) Handle(request eiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("PingRouter rev from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

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
	fmt.Println("Call HelloEarnthRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("HelloEarnthRouter recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Earnth Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnectionBegin 创建连接时执行
func DoConnectionBegin(conn eiface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ... ")
	//设置链接器
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnectionLost 连接断开的时候执行
func DoConnectionLost(conn eiface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}
	fmt.Println("DoConnectionLost is Called ... ")
}

func main() {
	s := enet.NewServer()
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloEarnthRouter{})
	s.Serve()
}