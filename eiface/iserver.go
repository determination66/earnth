package eiface

// IServer 定义Server接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 注册路由供客户端连接
	AddRouter(msgId uint32, router IRouter)
	// GetConnManager 得到链接管理
	GetConnManager() IConnManger

	//SetOnConnStart 设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//SetOnConnStop 设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//CallOnConnStart 调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//CallOnConnStop 调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
