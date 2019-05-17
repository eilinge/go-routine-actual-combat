package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	index = "http://pic.netbian.com"
)

// GetRealURL ...
func GetRealURL(url string) {
	// fmt.Printf("spider url: %s \n", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer res.Body.Close()

	var result string
	buffer := make([]byte, 1024)
	for {
		n, err := res.Body.Read(buffer)
		if err == io.EOF {
			// fmt.Println("res.Body.Read over!")
			break
		}
		result += string(buffer[:n])
	}

	re := regexp.MustCompile(`<a href="" id="img"><img src="(.*?)" data-pic=`)
	joyURL := re.FindAllStringSubmatch(result, -1)
	// fmt.Println("joyURL: ", joyURL)
	for _, eachURL := range joyURL {
		fmt.Println("eachURL: ", eachURL[1])
		// spiderURL = eachURL[1]
		// uploads/allimg/190422/105003-1555901403d3f8.jpg
		imgurl := index + eachURL[1]
		fmt.Println("imgurl: ", imgurl)

		imgName := strings.Split(strings.Split(imgurl, "/")[6], "-")[1]
		resp, err := http.Get(imgurl)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatal(resp.Status)
		}

		f, err := os.Create("C:\\Users\\wuchan4x\\Pictures\\carton\\" + imgName)
		if err != nil {
			log.Panic("文件创建失败")
		}
		io.Copy(f, resp.Body)
	}
}

// SpiderImage ...
func SpiderImage(url string) (err error) {
	// fmt.Printf("spider url: %s \n", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer res.Body.Close()

	var result string
	buffer := make([]byte, 1024)
	for {
		n, err := res.Body.Read(buffer)
		if err == io.EOF {
			// fmt.Println("res.Body.Read over!")
			break
		}
		result += string(buffer[:n])
	}

	// <a href="/tupian/24090.html" target="_blank"><img
	re := regexp.MustCompile(`<li><a href="(.*?)" target="_blank"><img src="/uploads/allimg/`)
	realURL := re.FindAllStringSubmatch(result, -1)
	for _, eachRealURL := range realURL {
		// fmt.Println("eachRealURL", eachRealURL[1])
		GetRealURL(index + eachRealURL[1])
	}

	return
}

// HTTPGet ...
func HTTPGet(i int, n chan int) {

	url := index + "/4kdongman/index_" + strconv.Itoa(i) + ".html"
	res, err := http.Get(url)
	// fmt.Println(url)
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer res.Body.Close()

	err = SpiderImage(url)
	if err != nil {
		fmt.Println("SpiderImage err: ", err)
		return
	}

	n <- i
}

// Dowork ...
func Dowork(start, end int) {
	// fmt.Printf("start: %d, end: %d\n", start, end)
	var n = make(chan int)
	t := time.NewTicker(10 * time.Second)
	for i := start; i <= end; i++ {

		// for {
		<-t.C

		go HTTPGet(i, n)

		fmt.Printf("the %d page spider done\n", <-n)
		// }
	}
	// 堵塞, 防止程序退出
	// for i := start; i <= end; i++ {
	// 	fmt.Printf("the %d page spider done\n", <-n)
	// }

}

func main() {
	fmt.Println("please input start page, this must > 1")
	var start, end int
	for {
		fmt.Scan(&start)
		if start < 1 {
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
