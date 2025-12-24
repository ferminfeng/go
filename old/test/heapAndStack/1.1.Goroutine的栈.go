package main

import (
	"fmt"
	"sync"
)

// 1.1 Goroutine的栈
// 在Go中，每个goroutine都有自己的栈，而不是每个线程。
func main() {
	fmt.Println("=== Goroutine栈演示 ===")

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 每个goroutine有自己的栈空间
			localVar := id * 100
			fmt.Printf("Goroutine %d: localVar=%d, 栈地址=%p\n",
				id, localVar, &localVar)

			// 调用函数，使用栈空间
			recursiveCall(id, 0)
		}(i)
	}

	wg.Wait()
}

func recursiveCall(goroutineID, depth int) {
	if depth >= 2 {
		return
	}

	// 每个递归调用都会在栈上分配新的变量
	localData := fmt.Sprintf("g%d-d%d", goroutineID, depth)
	fmt.Printf("  Goroutine %d 深度 %d: %s, 地址=%p\n",
		goroutineID, depth, localData, &localData)

	recursiveCall(goroutineID, depth+1)
}
