package main

import (
	"fmt"
	"reflect"
)

/*
反射reflection
	反射可大大提高程序的灵活性, 使得interface{} 有更大的发挥余地
	反射使用TypeOf 和 ValueOf函数从接口中获取目标对象信息
	反射会将匿名字段作为独立字段(匿名字段本质)
	想要利用反射修改对象状态, 前提是interface.data是settbale, 即pointer-interface
	通过反射可以"动态"调用方法
*/

type user struct {
	ID   int
	Name string
	Age  int
}

func (u user) hello() {
	fmt.Println("Hello world")
}

func main() {
	u := user{1, "eilinge", 17}
	info(u)
}

// 参数为空接口时, 可以传入任何类型
func info(o interface{}) {
	t := reflect.TypeOf(o)
	fmt.Println("Type:", t.Name())

	// 判断传入的类型是否正确
	if k := t.Kind(); k != reflect.Struct {
		// if k := t.Kind(); k != reflect.Map {
		fmt.Println("XXX")
		return
	}
	v := reflect.ValueOf(o)
	fmt.Println("Fields:")

	// 反射出类型的字段
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Printf("%6s: %v = %v\n", f.Name, f.Type, val)
	}

	// 反射出类型的方法
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf("%6s: %v\n", m.Name, m.Type)
	}
}
