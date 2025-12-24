package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// 创建两个 Person 结构体实例
	person1 := Person{Name: "Alice", Age: 30}
	// person2 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}

	// 使用 cmp.Equal 检查两个结构体是否相等
	equal := cmp.Equal(person1, person2)

	if equal {
		fmt.Println("The two people are equal.")
	} else {
		fmt.Println("The two people are not equal.")
	}
}
