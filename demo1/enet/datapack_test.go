package enet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {

	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				fmt.Println("server accept err:", err)
				return
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err:", err)
						return
					}

					//将headData字节流 拆包到msg中
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						//msg 是有data数据的，需要再次读取data数据
						msg, ok := msgHead.(*Message)
						if !ok {
							fmt.Println("msgHead assert failed!")
						}
						msg.Data = make([]byte, msg.GetDataLen())

						//根据dataLen从io中读取字节流
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}

				}
			}(conn)

		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dail err,", err)
		return
	}

	dp := NewDataPack()

	//创建粘包情况，并发送
	msgs := []*Message{
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

	conn.Write(sendData)

	//客户端阻塞
	select {}

}
