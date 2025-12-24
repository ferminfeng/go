package main

import (
	"fmt"
)

// 编译时分析: go build -gcflags="-m" main.go
// 运行时分析: go run -gcflags="-m" main.go

func main() {
	fmt.Println("=== 逃逸分析实战 ===")

	example1() // 栈分配
	example2() // 堆分配 - 指针逃逸
	example3() // 堆分配 - 接口逃逸
	example4() // 堆分配 - 闭包逃逸
	example5() // 堆分配 - 大小逃逸
}

// 示例1: 栈分配 - 变量没有逃逸
func example1() {
	x := 100
	y := 200
	sum := x + y
	fmt.Printf("栈分配示例: %d\n", sum)
	// x, y, sum 都在栈上
}

// 示例2: 堆分配 - 指针逃逸
func example2() *int {
	v := 42
	return &v // v逃逸到堆上，因为它在函数返回后还需要存在
}

// 示例3: 堆分配 - 接口逃逸
func example3() {
	data := "hello" // data逃逸到堆上，因为被接口使用
	var anything interface{} = data
	fmt.Println(anything)
}

// 示例4: 堆分配 - 闭包逃逸
func example4() func() int {
	counter := 0 // counter逃逸到堆上，因为被闭包捕获
	return func() int {
		counter++
		return counter
	}
}

// 示例5: 堆分配 - 大小逃逸
func example5() {
	small := [1]int{1}    // 可能在栈上
	large := [10000]int{} // 可能在堆上，因为太大

	fmt.Printf("小数组: %p, 大数组: %p\n", &small, &large)
}
