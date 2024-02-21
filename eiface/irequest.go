package eiface

// IRequest 封装客户端请求信息和请求数据
type IRequest interface {
	GetConnection() IConnection

	GetData() []byte

	GetMsgID() uint32
}
