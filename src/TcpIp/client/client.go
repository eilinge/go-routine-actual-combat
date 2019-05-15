package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("net.Dial err: ", err)
		return
	}
	defer conn.Close()

	// 接收服务器回复的数据, 开启新任务(go routine)
	buffer := make([]byte, 1024)
	go func() {
		for {
			n, err := conn.Read(buffer) // 接收服务器的请求(request)
			if err != nil {
				log.Fatal("conn.Read err: ", err)
				// return
				break
			}

			fmt.Printf("reponse = %+v", string(buffer[:n])) // 显示response
		}
	}()

	// 从键盘输入内容, 给服务器发送request
	str := make([]byte, 1024)
	// go func() { // 非堵塞
	for { // 堵塞
		n, err := os.Stdin.Read(str)
		if err != nil {
			log.Fatal("os.Stdin.Read err: ", err)
			// return
			break
		}
		conn.Write(str[:n]) // response:[]byte
	}
	// }()
}
