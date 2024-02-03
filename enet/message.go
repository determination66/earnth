package enet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(msgId uint32) {
	m.Id = msgId
}

func (m *Message) setData(msgData []byte) {
	m.Data = msgData
}

func (m *Message) setDataLen(msgLen uint32) {
	m.DataLen = msgLen
}
