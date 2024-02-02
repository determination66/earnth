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

	Version        string
	MaxConn        int
	MaxPackageSize uint32
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("config/earnth.json")
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
		Version:        "v0.4",
		TcpPort:        8888,
		Host:           "0:0:0:0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()

}
