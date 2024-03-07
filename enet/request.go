package enet

import "earnth/eiface"

type Request struct {
	//已经和客户端建立好的连接
	conn eiface.IConnection

	//客户端请求的数据
	msg eiface.IMessage
}

func (r *Request) GetConnection() eiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//	func (r *Request) GetMsgData() uint32 {
//		return r.msg.GetMsgId()
//	}
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
