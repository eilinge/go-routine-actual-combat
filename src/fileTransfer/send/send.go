package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// SendFile ...
func SendFile(path string, conn net.Conn) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("os.Open err: ", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 2048)

	for {
		n, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Println("file read over")
			} else {
				fmt.Println("file.Read err: ", err)
			}
			return
		}
		// 将文件全部发送
		conn.Write(buffer[:n])
	}
}

func main() {

	// fmt.Println("请输入需要传输的文件: ")
	var path string

	list := os.Args
	if len(list) != 2 {
		fmt.Printf("useage: srcfile:%s \n", list)
		return
	}
	// list[0] 该执行文件
	path = list[1]

	// fmt.Scan(&path)

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal("os.Stat err: ", err)
		return
	}
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("net.Dial err: ", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(info.Name()))
	if err != nil {
		log.Fatal("conn.Write err: ", err)
		return
	}

	buffer := make([]byte, 1024)
	// for {
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatal("conn.Read err: ", err)
		return
	}

	if "ok" == string(buffer[:2]) {
		// fmt.Println("rending send file ...")
		// go SendFile(path, conn)
		SendFile(path, conn)
	}
	// }
}
