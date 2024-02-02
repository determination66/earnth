package enet

import "earnth/eiface"

type Request struct {
	//已经和客户端建立好的连接
	conn eiface.IConnection

	//客户端请求的数据
	data []byte
}

func (r *Request) GetConnect() eiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
