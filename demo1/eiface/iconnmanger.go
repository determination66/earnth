package eiface

type IConnManger interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connId uint32) (IConnection, error)
	Len() int
	CLearConn()
}
