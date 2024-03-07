package enet

import (
	"earnth/eiface"
	"earnth/utils"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IpVersion string
	Ip        string
	Port      int
	//当前的Server的消息管理模块
	msgHandler eiface.IMsgHandler
}

// Start 启动Server
func (s *Server) Start() {
	go func() {
		fmt.Printf("[Start] Server Listener at IP: %s,Port: %d, is starting...\n", s.Ip, s.Port)
		//1.基本服务器开发，获取Tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Printf("resolve tcp addr error:%v", err)
			return
		}
		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Printf("listenTCP error:%v...", err)
			return
		}
		fmt.Printf("start earnth server success:%s...\n", s.Name)

		var cid uint32
		cid = 0

		//3.阻塞等待客户端连接，处理客户端业务(读写)
		for {
			//如果客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			dealConn := NewConnection(conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}

	}()

}

func (s *Server) Stop() {
	//todo 停止服务器，状态资源，已经开辟的信息的停止
}

func (s *Server) Serve() {
	s.Start()
	// 后续可以做一些其他业务
	select {}
}

func (s *Server) AddRouter(msgId uint32, router eiface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add router success! msgId = ", msgId)
}

// NewServer 初始化Server模块方法
func NewServer() eiface.IServer {

	return &Server{
		Name:       utils.GlobalObject.Name,
		IpVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
	}

}
