package main

import (
	"earnth/enet"
	"fmt"
	"net"
	"time"
)

// 客户端
func main() {
	fmt.Println("client start...")

	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net dail err:", err)
		return
	}
	dp := enet.NewDataPack()

	for {
		msgs := []*enet.Message{
			{Id: 0, DataLen: 5, Data: []byte{'h', 'e', 'l', 'l', 'o'}},
			{Id: 1, DataLen: 5, Data: []byte{'w', 'o', 'r', 'l', 'd'}},
		}
		var sendData []byte
		for _, v := range msgs {
			data, err := dp.Pack(v)
			if err != nil {
				fmt.Println("pack err:", err)
				return
			}
			sendData = append(sendData, data...)
		}

		_, err := conn.Write(sendData)
		//_, err := conn.Write([]byte("hello,world!"))
		if err != nil {
			fmt.Println("conn write err:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn read err:", err)
			return
		}
		fmt.Println("server call back:", string(buf), ", cnt:", cnt)
		time.Sleep(time.Second)
	}

}
