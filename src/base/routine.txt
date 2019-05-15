package main

import (
	"fmt"
	"runtime"
	"time"
)

func test() {
	defer fmt.Println("bbbbbb")
	// return
	runtime.Goexit() // 终止所在协程
	fmt.Println("dddddd")
}

func main01() {
	go func() {
		for i := 0; i < 1; i++ {
			fmt.Println("aaaaaa")
			test()
			fmt.Println("cccccc")
		}
	}()

	for j := 0; j < 3; j++ {
		runtime.Gosched() // 让出时间片, 让别的协程执行, 执行完毕, 再回来执行
		fmt.Println("hello")
	}
}

func main02() {
	runtime.GOMAXPROCS(4) // 设置执行的核数
	for {
		go func() {
			fmt.Print(0)
		}()
		fmt.Print(1)
	}
}

// Printor ...
func Printor(str string) {
	for _, data := range str {
		fmt.Printf("%c", data)
		time.Sleep(time.Millisecond)
	}
}

func person1(ch chan<- byte) {
	Printor("hello")
	ch <- 'a'
}

func person2(ch <-chan byte) {
	<-ch
	Printor("world")
}

func main() {
	ch := make(chan byte)
	go person1(ch)
	go person2(ch)

	time.Sleep(time.Second)
}
