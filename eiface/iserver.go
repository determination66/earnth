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
	AddRouter()
}
