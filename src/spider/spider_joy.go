package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

// SpiderJoy ...
func SpiderJoy(url string) (spiderURL string, err error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer res.Body.Close()

	var result string
	buffer := make([]byte, 1024)
	for {
		n, _ := res.Body.Read(buffer)
		if n == 0 {
			break
		}
		result += string(buffer[:n])
	}

	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))" target="_blank">`)
	joyURL := re.FindAllStringSubmatch(result, -1)
	// fmt.Println("joyURL: ", joyURL)
	for _, eachURL := range joyURL {
		// fmt.Println("eachURL: ", eachURL)
		spiderURL = eachURL[1]
	}
	return
}

// GetTitleContent title, content, err := GetTitleContent(result)
func GetTitleContent(result string) (title, content string, err error) {
	return
}

// HTTPGet ...
// func HTTPGet(i int, n chan int) {
func HTTPGet(i int) {

	url := "https://www.pengfu.com/xiaohua_" + strconv.Itoa(i) + ".html"
	res, err := http.Get(url)
	// fmt.Println(url)
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer res.Body.Close()

	_, err = SpiderJoy(url)
	if err != nil {
		fmt.Println("SpiderJoy err: ", err)
		return
	}
	// fmt.Println("result: ", result)
	// title, content, err := GetTitleContent(result)
	// fmt.Printf("title: %s", title)
	// fmt.Printf("content: %s", content)
	// n <- i
}

// Dowork ...
func Dowork(start, end int) {
	// fmt.Printf("start: %d, end: %d\n", start, end)
	// var n = make(chan int)
	for i := start; i <= end; i++ {
		HTTPGet(i)
	}
	// // 堵塞, 防止程序退出
	// for i := start; i <= end; i++ {
	// 	<-n
	// }

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
