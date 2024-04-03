package main

import (
	"earnth/enet"
	"fmt"
	"io"
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
			//{Id: 1, DataLen: 5, Data: []byte{'w', 'o', 'r', 'l', 'd'}},
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

		//重新获取服务端发送来的内容
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*enet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Client receive Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}

}
