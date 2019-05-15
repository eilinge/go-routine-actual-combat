package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Recvfile ...
func Recvfile(fileName string, conn net.Conn) {
	df, _ := os.Create(fileName)
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, "connect successfully")

	defer df.Close()

	buf := make([]byte, 4*1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("receive over!")
			} else {
				fmt.Println("err = ", err)
			}
			break
		}
		df.Write(buf[:n])
	}
	conn.Write([]byte("receive over!"))
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("net.Listen err: ", err)
		return
	}

	defer listener.Close()

	// for {
	conn, err := listener.Accept()
	defer conn.Close()
	if err != nil {
		log.Fatal("listener.Accept err: ", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal("conn.Read err: ", err)
		return
	}

	fileName := string(buffer[:n])

	conn.Write([]byte("ok"))

	Recvfile(fileName, conn)
	// }
}
