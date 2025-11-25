package main

import "fmt"

// 1.2 栈的动态增长
// Go的栈不是固定大小的，而是可以动态增长的：
func main() {
	demonstrateStackGrowth()
}

func demonstrateStackGrowth() {
	fmt.Println("=== 栈动态增长演示 ===")

	// 初始状态
	var initial [128]byte
	fmt.Printf("初始栈变量地址: %p\n", &initial)

	// 深度递归，触发栈增长
	deepRecursion(0)
}

func deepRecursion(depth int) {
	var localArray [256]byte // 占用栈空间

	if depth%100 == 0 {
		fmt.Printf("深度 %d: 栈变量地址=%p\n", depth, &localArray)
	}

	if depth < 1000 {
		deepRecursion(depth + 1)
	}
}
