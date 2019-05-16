package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// HTTPGet ...
func HTTPGet(i int, n chan int) {

	url := "https://tieba.baidu.com/f?kw=golang&ie=utf-8&pn=" + strconv.Itoa((i-1)*50)
	res, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer res.Body.Close()

	result := make([]byte, 1024)
	buffer := make([]byte, 1024)
	for {
		n, _ := res.Body.Read(buffer)
		if n == 0 {
			break
		}
		result = append(result, buffer[:n]...)
	}
	dstfile := strconv.Itoa(i) + ".html"
	df, _ := os.Create(dstfile)
	defer df.Close()
	df.Write(result)
	fmt.Printf("spider page: %d over!\n", i)
	n <- i
}

// Dowork ...
func Dowork(start, end int) {
	// fmt.Printf("start: %d, end: %d\n", start, end)
	var n = make(chan int)
	for i := start; i <= end; i++ {
		go HTTPGet(i, n)
	}
	// 堵塞, 防止程序退出
	for i := start; i <= end; i++ {
		<-n
	}

}

func main() {
	fmt.Println("please input start page, this must >= 1")
	var start, end int
	for {
		fmt.Scan(&start)
		if start <= 0 {
			fmt.Println("input start error, please input agin")
			continue
		}
		fmt.Println("please input end page, this must >= start")
		fmt.Scan(&end)
		if end < start {
			fmt.Println("input end error, please input agin")
			continue
		}

		Dowork(start, end)
		return
	}

}
