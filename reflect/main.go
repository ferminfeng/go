package main

import (
	"fmt"
	"reflect"
)

type Cal struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

func (c Cal) GetSub(name string) {
	fmt.Printf("%v完成了减法运行，%v-%v=%v\n", name, c.Num1, c.Num2, c.Num1-c.Num2)
}

func reflectCal(a interface{}) {
	// 使用反射遍历结构体所有的字段信息
	refType := reflect.TypeOf(a)
	value := reflect.ValueOf(a)
	if value.Kind() != reflect.Struct {
		// fmt.Println("进来了！")
		return
	}
	// 获取字段个数
	num := value.NumField()
	// 遍历字段信息
	for i := 0; i < num; i++ {
		fmt.Printf("字段个数是：%v,字段名称是：%v\n", num, refType.Field(i))
	}

	// 查看该value有几个方法
	numOfMethod := value.NumMethod()
	fmt.Printf("\nvalue有%v方法", numOfMethod)

	// //定义一个反射切片
	// var refValue []reflect.Value
	// //设置cal结构体字段的值并添加到反射切片
	// refValue = append(refValue, reflect.ValueOf(8))
	// refValue = append(refValue, reflect.ValueOf(5))
	//
	// sub := value.Method(0).Call(refValue)
	// fmt.Printf("sub的值是：%v", sub)

}

func main() {
	// var cal Cal = Cal{}

	reflectCal(Cal{})
}
