package eiface

type IMsgHandler interface {
	//DoMsgHandler 马上以非阻塞方式处理消息
	DoMsgHandler(request IRequest)
	//AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
}
