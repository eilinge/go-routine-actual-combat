package main

import (
	"fmt"
	"regexp"
)

func main01() {
	buffer := "abc a1c abac a8c aaa"

	// reg1 := regexp.MustCompile(`a.c`)
	// reg1 := regexp.MustCompile(`a[0-9]c`)
	// reg1 := regexp.MustCompile(`a\wc`)
	// reg1 := regexp.MustCompile(`a\dc`)
	reg1 := regexp.MustCompile(`a\Wc`)
	if reg1 == nil {
		panic("regxep match failed")
	}

	// n int: 匹配次数
	// < 0: 全部匹配
	// n > 0: 匹配n次
	result := reg1.FindAllStringSubmatch(buffer, -1)

	fmt.Println("result = ", result) // [[abc] [a1c] [a8c]]
}

func main() {
	buff := `<body id="t" style="font-family: 'Segoe UI', Tahoma, sans-serif; font-size: 75%" jstcache="0" class="neterror">
  <div id="main-frame-error" class="interstitial-wrapper" jstcache="0">
    <div id="main-content" jstcache="0">
      <div class="icon icon-generic" jseval="updateIconClass(this.classList, iconClass)" alt="" jstcache="1"></div>
      <div id="main-message" jstcache="0">
		<h1 jstcache="0">
		  <div>哈哈</div>
		  <div>我最帅
		  低调如我
		  2333
		  </div>
		  <div>喔喔</div>
          <span jsselect="heading" jsvalues=".innerHTML:msg" jstcache="10">This site can’t be reached</span>
          <a id="error-information-button" class="hidden" onclick="toggleErrorInformationPopup();" jstcache="0"></a>
        </h1>
		<p jsselect="summary" jsvalues=".innerHTML:msg" `

	// ?s: 匹配换行符(\n)
	reg1 := regexp.MustCompile(`<div>(?s:(.*?))</div>`)
	if reg1 == nil {
		panic("regxep match failed")
	}

	// n int: 匹配次数
	// < 0: 全部匹配
	// n > 0: 匹配n次
	result := reg1.FindAllStringSubmatch(buff, -1)

	// fmt.Println("result = ", result) // [[abc] [a1c] [a8c]]
	/*
			[[<div>哈哈</div> 哈哈] [<div>我最帅
		                  低调如我
		                  2333
		                  </div> 我最帅
		                  低调如我
		                  2333
						  ] [<div>喔喔</div> 喔喔]]
	*/
	for _, v := range result {
		fmt.Println("v[0] = ", v[0]) // v[0] =  <div>哈哈</div>
		fmt.Println("v[1] = ", v[1]) // v[1] =  哈哈
	}
}
