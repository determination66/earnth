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
	//当前的Server的链接管理器
	ConnManger eiface.IConnManger

	//新增两个hook函数原型
	//该Server的连接创建时Hook函数
	OnConnStart func(conn eiface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn eiface.IConnection)
}

//对Server链接的创建和断开后进行个性化的处理

func (s *Server) SetOnConnStart(f func(eiface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(eiface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(conn eiface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn eiface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

// Start 启动Server
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.Ip, s.Port)
	fmt.Printf("[EARNTH] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		// 启动worker工作池
		s.msgHandler.StartWorkerPool()
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

			if s.ConnManger.Len() > utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}

	}()

}

func (s *Server) Stop() {
	//todo 停止服务器，状态资源，已经开辟的信息的停止
	fmt.Println("[STOP] Earnth server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnManger.CLearConn()
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

func (s *Server) GetConnManager() eiface.IConnManger {
	return s.ConnManger
}

// NewServer 初始化Server模块方法
func NewServer() eiface.IServer {

	return &Server{
		Name:       utils.GlobalObject.Name,
		IpVersion:  "tcp4",
		Ip:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
		ConnManger: NewConnManager(),
	}

}
