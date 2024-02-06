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
	//创建一个封包对象 dp
	dp := enet.NewDataPack()

	//封装一个msg1包
	msg1 := &enet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &enet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)

	//向服务器端写数据
	conn.Write(sendData1)

	//客户端阻塞
	select {}
	//for {
	//	_, err := conn.Write([]byte("hello,world!"))
	//	if err != nil {
	//		fmt.Println("conn write err:", err)
	//		return
	//	}
	//	buf := make([]byte, 512)
	//	cnt, err := conn.Read(buf)
	//	if err != nil {
	//		fmt.Println("conn read err:", err)
	//		return
	//	}
	//	fmt.Println("server call back:", string(buf), ", cnt:", cnt)
	//	time.Sleep(time.Second)
	//}

}
