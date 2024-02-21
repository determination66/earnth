package main

import (
	"earnth/eiface"
	"earnth/enet"
	"fmt"
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
	fmt.Println("********", req.GetData())
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
	s := enet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
