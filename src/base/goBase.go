package main

import "fmt"

const (
	a       = iota
	i, j, k = iota, iota, iota
	m       = iota
)

func main01() {
	fmt.Printf("a type is: %T\n", a)
	fmt.Printf("a = %d, i = %d, j = %d, k = %d, m = %d\n", a, i, j, k, m) // a = 0, i = 1, j = 1, k = 1, m = 2

	// 浮点型
	var f1 float32
	f1 = 3.14
	fmt.Println("f1 = ", f1)

	f2 := 3.14
	fmt.Printf("f2 type:%T, value:%v\n", f2, f2) // 自动推导式默认:float64

	var ch byte
	ch = 97
	// fmt.Println("ch = ", ch)
	// 格式化输出, %c以字符方式打印, %d以整型方式打印
	fmt.Printf("char type:%c , int type:%d\n", ch, ch) // char type:a , int type:97

	// 字符, 单引号
	ch = 'a'
	fmt.Printf("%c, %d\n", ch, ch) // a, 97

	//大小写转换, 相差32, 小写大
	fmt.Printf("upper:%d, lower:%d\n", 'A', 'a') // upper:65, lower:97
	fmt.Printf("upper convert lower:%c\n", 'A'+32)
	fmt.Printf("lower convert upper:%c\n", 'a'-32)

	var str string
	//1. 双引号
	//2. 字符串有1个或多个字符组成
	//3. 字符串都是隐藏一个结束符, '\0'
	str = "a"
	fmt.Println("str = ", str)
	str = "hello world"
	// 只想操作字符串的某个字符, 从0开始操作
	fmt.Printf("str[0] = %c, str[1] = %c\n", str[0], str[1])

	// var t complex64
	// t = 2.1 + 3.14i
	t2 := 2.1 + 3.14i // t2 type is complex128
	fmt.Printf("t2 type is %T\n", t2)

	// 通过内建函数, 取实部和虚部
	// real(t2) =  2.1 , imag(t2) =  3.14
	fmt.Println("real(t2) = ", real(t2), ", imag(t2) = ", imag(t2))

	var d int
	fmt.Printf("请输入变量a: ")
	// fmt.Scanf("%d", &d)
	fmt.Scan(&d)
	fmt.Println("a = ", d)

	var flag bool
	flag = true
	fmt.Printf("flag = %t\n", flag)

	//bool 类型不能转换为int
	// fmt.Printf("flag = %d", flag)

	// 0就是假, 非0就是真
	// 整型也不能转换为bool
	// flag = bool(1)

	var ch1 byte
	ch1 = 'a'
	var t int
	t = int(ch1)
	fmt.Println("t = ", t)
}

func main02() {
	// var num int
	// fmt.Printf("请按下楼层: ")
	// fmt.Scan(&num)
	// switch num {
	// case 1:
	// 	fmt.Println("按下的是1楼")
	// case 2:
	// 	fmt.Println("按下的是2楼")
	// 	// 默认:break, 直接退出循环体
	// case 3:
	// 	fmt.Println("按下的是3楼")
	// 	fallthrough // 不跳出switch语句, 后面无条件执行
	// case 4:
	// 	fmt.Println("按下的是4楼")
	// 	fallthrough
	// default:
	// 	fmt.Println("是其他楼层")
	// }

	switch num := 3; num {
	case 1:
		fmt.Println("按下的是1楼")
	case 2:
		fmt.Println("按下的是2楼")
		// 默认:break, 直接退出循环体
	case 3, 4, 5, 6:
		fmt.Println("按下的是yyyy楼")
		fallthrough // 不跳出switch语句, 后面无条件执行
	case 7:
		fmt.Println("按下的是4楼")
		fallthrough
	default:
		fmt.Println("是其他楼层")
	}

	score := 85
	switch {
	case score > 90:
		fmt.Println("优秀")
	case score > 80:
		fmt.Println("良好")
	case score > 70:
		fmt.Println("一般")
	default:
		fmt.Println("其他")
	}
}

func main03() {
	str := "abcefg"
	for i := 0; i < len(str); i++ {
		fmt.Printf("str[%d]=%c\n", i, str[i])
	}

	// 迭代打印每个元素, 默认返回2个值: key, value
	for i, data := range str {
		fmt.Printf("str[%d]=%c\n", i, data)
	}

	// 第2个返回值, 默认丢弃, 返回元素的位置(key)
	for i := range str {
		if i == 2 {
			// str[0]=a
			// str[1]=b
			// break // 跳出循环, 跳出最近的内循环

			// str[0]=a
			// str[1]=b
			// str[3]=e
			// str[4]=f
			// str[5]=g
			continue // 跳过本次循环, 下一次继续
		}
		fmt.Printf("str[%d]=%c\n", i, str[i])
	}

	// for i, _ := range str {
	// 	fmt.Printf("str[%d]=%c\n", i, str[i])
	// }
}

func main04() {
	// break // break is not in a loop, switch, or select
	// continue // continue is not in a loop

	fmt.Println("12312312")

	goto ERR // 可用于任何地方, 但是不允许跨函数使用, 调整执行位置,
	// fmt.Println("23333333") // unreachable code

ERR:
	fmt.Println("44444444")
}

func test(a int) {
	if a == 1 {
		fmt.Println("value of a is: ", a)
		return
	}
	test(a - 1)
	// 退出函数递归之后, 才逆顺序执行
	fmt.Println("a value:", a)

	// value of a is:  1
	// a value: 2
	// a value: 3
	// a value: 4
	// this func called end
}
func main() {
	test(4)
	fmt.Println("this func called end")
}
