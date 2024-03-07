package enet

import (
	"earnth/eiface"
	"fmt"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]eiface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]eiface.IRouter),
	}
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (mh *MsgHandler) DoMsgHandler(request eiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}
	//执行对应方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

}

// AddRouter 为消息添加具体的逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router eiface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}
