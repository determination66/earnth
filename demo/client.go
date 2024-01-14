package main

import (
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
	for {
		_, err := conn.Write([]byte("hello,world!"))
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
