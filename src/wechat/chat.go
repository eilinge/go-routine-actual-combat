package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Client ...
type Client struct {
	C    chan string // 用户发送数据的管道(中转channel), 发送消息给每一个
	Name string      // 用户名
	Addr string      // 地址
}

// 在线人数
var onlineMap map[string]Client

// 广播通讯(广播channel)
var message = make(chan string)

// MessageToEachClient 转发消息, 只要有消息来了, 遍历map, 给map每个成员都发送此消息
func MessageToEachClient() {
	// 给map分配空间: onlineMap not nil
	onlineMap = make(map[string]Client)
	for {
		msg := <-message

		// 遍历map, 给map每个成员都发送此消息
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}
}

// WriteMsgToClient ...
func WriteMsgToClient(cli Client, conn net.Conn) {
	for msg := range cli.C {
		conn.Write([]byte(msg + "\n"))
	}
}

// MakeMsg ...
func MakeMsg(cli Client, msg string) (buf string) {
	buf = "[" + cli.Addr + "]" + cli.Name + ": " + msg
	return buf
}

// HandleConn ...
func HandleConn(conn net.Conn) { // 处理用户连接
	defer conn.Close()
	cliAddr := conn.RemoteAddr().String()
	cli := Client{make(chan string), cliAddr, cliAddr}

	// 添加新成员
	onlineMap[cliAddr] = cli

	// send message to each client
	go WriteMsgToClient(cli, conn)

	isquit := make(chan bool)
	isData := make(chan bool)
	// each client receive and send broadcast message
	go func() {
		buffer := make([]byte, 1024)

		for {
			n, err := conn.Read(buffer)
			if n == 0 { // 对方断开, 或出问题
				isquit <- true
				fmt.Println("Conn.Read err: ", err)
				break
			}
			msg := string(buffer[:n-1])

			if len(msg) == 3 && msg == "who" {
				for _, tmp := range onlineMap {
					conn.Write([]byte(tmp.Name + "\n"))
				}
			} else if strings.Split(msg, "|")[0] == "rename" {
				cli.Name = strings.Split(msg, "|")[1]
				onlineMap[cliAddr] = cli // 修改map
				conn.Write([]byte("rename ok!\n"))
			} else {
				message <- MakeMsg(cli, msg)
			}
			isData <- true // 输入数据, 用户阻塞
		}
	}()

	// 广播某个在线
	message <- MakeMsg(cli, "login")

	for { // 防止广播一次之后, 连接断开 use of closed network connection
		select {
		case <-isquit: // 用户退出
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, "logout")
			break
		case <-isData: // 有数据输入, 堵塞, 防止程序退出
		case <-time.After(5 * time.Millisecond): // 延迟处理
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, "time out")
			break
		}
	}
}

// HandleConn -> MessageToEachClient -> WriteMsgToClient
// 对消息进行分类处理
// 开启协程: 1.分别处理每一个Conn; 2.进程间通过channel通信
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	// 新开一个协程, 转发消息, 只要有消息来了, 遍历map, 给map每个成员都发送此消息
	go MessageToEachClient()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			continue
		}

		defer conn.Close()
		go HandleConn(conn)
	}
}
