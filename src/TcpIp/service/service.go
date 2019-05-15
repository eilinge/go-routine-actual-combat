package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// HandleConn ...
func HandleConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, "connect successfully")

	buffer := make([]byte, 1024) // 缓存区

	for {
		n, err := conn.Read(buffer) // 接收服务器的请求(request)
		if err != nil {
			log.Fatal("conn.Read err: ", err)
			// return
			break
		}

		fmt.Printf("buffer = %+v", string(buffer[:n]))
		// fmt.Println(len(buffer))
		if string(buffer[:4]) == "exit" {
			fmt.Println("conn close")
			break
		}
		conn.Write([]byte(strings.ToUpper(string(buffer[:n])))) // 返回响应(response)给client: []byte
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("net.Listen err: ", err)
		return
	}

	defer listener.Close()
	// go func() {
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal("listener.Accept err: ", err)
			return
		}
		// 并发处理request
		go HandleConn(conn)
	}
	// }()
}
