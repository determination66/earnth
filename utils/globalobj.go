package utils

import (
	"earnth/eiface"
	"encoding/json"
	"os"
)

type GlobalObj struct {
	//Server配置
	TcpServer eiface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version string
	//最大链接数
	MaxConn int
	//最大
	MaxPacketSize uint32
	//worker的工作数量
	WorkerPoolSize uint32
	//每个Task的最大长度
	MaxWorkerTaskLen uint32
	// 缓冲管道数目
	MaxMsgChanLen uint32

	ConfigFilePath string
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile(g.ConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

// 定义全局对外GlobalObj

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "earthServerApp",
		Version:        "v0.9",
		TcpPort:        8888,
		Host:           "0:0:0:0",
		MaxConn:        1000,
		MaxPacketSize:  4096,
		ConfigFilePath: "config/earnth.json",
	}
	GlobalObject.Reload()
}
