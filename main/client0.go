package main

import (
	"earnth/enet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client0 Test!!!")
	//time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	dp := enet.NewDataPack()
	for {
		msg, _ := dp.Pack(enet.NewMsgPackage(0, []byte("Client0 Package")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err:", err)
			return
		}

		//先读取流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head err:", err)
			return
		}

		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*enet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack err:", err)
				return
			}
			fmt.Println("==> client0 Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}
		time.Sleep(1 * time.Second)

	}

}
