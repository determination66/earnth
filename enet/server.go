package enet

import (
	"earnth/eiface"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IpVersion string
	Ip        string
	Port      int
}

func CallBack(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("Callback is called!")
	_, err := conn.Write(data[:cnt])
	if err != nil {
		return errors.New("CallBack to client err")
	}
	return nil
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
		fmt.Printf("start earnth server success:%s...", s.Name)

		var cid uint32
		cid = 0

		//NewConnection()

		//3.阻塞等待客户端连接，处理客户端业务(读写)
		for {
			//如果客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			dealConn := NewConnection(conn, cid, CallBack)
			cid++

			go dealConn.Start()

			////客户端已经建立连接，conn,最基本最大512回写业务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("Read err:", err)
			//			continue
			//		}
			//		fmt.Println("client send:", string(buf))
			//		// 回写功能
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back err:", err)
			//			continue
			//		}
			//	}
			//}()

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

// NewServer 初始化Server模块方法
func NewServer(name string) eiface.IServer {
	return &Server{
		Name:      name,
		IpVersion: "tcp4",
		Ip:        "0.0.0.0",
		Port:      8888,
	}

}
