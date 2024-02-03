package eiface

// IMessage 封装消息
type IMessage interface {
	GetDataLen() uint32
	GetMsgId() uint32
	GetData() []byte

	SetMsgId(uint322 uint32)
	setData([]byte)
	setDataLen(uint322 uint32)
}
