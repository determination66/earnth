package eiface

type IDataPack interface {
	// GetHeadLen 获取包头的长度
	GetHeadLen() uint32

	// Pack 封包操作
	Pack(msg IMessage) ([]byte, error)

	// UnPack 拆包操作
	UnPack([]byte) (IMessage, error)
}
