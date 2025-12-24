package main

import (
	"fmt"
	"testing"
	"time"
)

type SmallStruct struct {
	a, b, c int32
}

type LargeStruct struct {
	data [1000]int64
}

// 基准测试：栈分配性能
func BenchmarkStackAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 小结构体，在栈上分配
		var local SmallStruct
		local.a = int32(i)
		local.b = local.a * 2
		local.c = local.b + local.a
		_ = local
	}
}

// 基准测试：堆分配性能
func BenchmarkHeapAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 通过指针使结构体在堆上分配
		heapObj := &SmallStruct{
			a: int32(i),
			b: int32(i * 2),
			c: int32(i * 3),
		}
		_ = heapObj
	}
}

// 大对象测试
func BenchmarkLargeStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 大数组，可能在堆上
		var large LargeStruct
		for j := range large.data {
			large.data[j] = int64(j)
		}
		_ = large
	}
}

func BenchmarkLargeHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 大对象在堆上分配
		large := &LargeStruct{}
		for j := range large.data {
			large.data[j] = int64(j)
		}
		_ = large
	}
}

func main() {
	// 运行简单的性能测试
	fmt.Println("性能测试结果:")

	// 测试小对象
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		var obj SmallStruct
		obj.a = int32(i)
		_ = obj
	}
	stackTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 1000000; i++ {
		obj := &SmallStruct{a: int32(i)}
		_ = obj
	}
	heapTime := time.Since(start)

	fmt.Printf("小对象 - 栈: %v, 堆: %v, 堆比栈慢: %.2fx\n",
		stackTime, heapTime, float64(heapTime)/float64(stackTime))

	// 测试大对象
	start = time.Now()
	for i := 0; i < 1000; i++ {
		var large LargeStruct
		_ = large
	}
	largeStackTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		large := &LargeStruct{}
		_ = large
	}
	largeHeapTime := time.Since(start)

	fmt.Printf("大对象 - 栈: %v, 堆: %v, 差异: %.2fx\n",
		largeStackTime, largeHeapTime,
		float64(largeHeapTime)/float64(largeStackTime))
}
