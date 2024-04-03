package utils

import (
	"earnth/demo1/eiface"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type GlobalObj struct {
	//Server配置
	TcpServer        eiface.IServer
	Host             string `yaml:"Host"`             // 主机地址
	TcpPort          int    `yaml:"TcpPort"`          // TCP端口
	Name             string `yaml:"Name"`             // 名称
	Version          string `yaml:"Version"`          // 版本
	MaxConn          int    `yaml:"MaxConn"`          // 最大连接数
	MaxPacketSize    uint32 `yaml:"MaxPacketSize"`    // 最大数据包大小
	WorkerPoolSize   uint32 `yaml:"WorkerPoolSize"`   // Worker池大小
	MaxWorkerTaskLen uint32 `yaml:"MaxWorkerTaskLen"` // 每个Worker的最大任务长度
	MaxMsgChanLen    uint32 `yaml:"MaxMsgChanLen"`    // 最大消息通道长度
	ConfigFilePath   string // 配置文件路径
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile(g.ConfigFilePath)
	if err != nil {
		panic(err)
	}
	// 打印解析前的数据
	fmt.Println("Raw earnth.yaml:", string(data))

	err = yaml.Unmarshal(data, &GlobalObject)

	if err != nil {
		panic(err)
	}

	// 打印解析后的数据
	fmt.Println("Parsed GlobalObject:", GlobalObject)
}

// 定义全局对外GlobalObj

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "earthServerApp",
		Version:        "v0.9",
		TcpPort:        8888,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPacketSize:  4096,
		ConfigFilePath: "demo1/config/earnth.yaml",
	}
	GlobalObject.Reload()
}
