package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	runtime.GOMAXPROCS(4)

	var w sync.WaitGroup
	count := int32(0)
	w.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				// 多核情况下 A核修改 count 的时候，由于 CPU 缓存的存在，B核读到的 count 值可能不是最新的值
				//count++
				atomic.AddInt32(&count, 1)
			}
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println(count)
}
