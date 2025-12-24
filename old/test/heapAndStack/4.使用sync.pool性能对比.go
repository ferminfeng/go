package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Message struct {
	ID      int
	Content string
	Headers map[string]string
	Data    []byte
}

// 无池版本
func createMessageWithoutPool(id int) *Message {
	return &Message{
		ID:      id,
		Content: "Hello World",
		Headers: make(map[string]string),
		Data:    make([]byte, 1024),
	}
}

// 有池版本
var messagePool = sync.Pool{
	New: func() interface{} {
		return &Message{
			Headers: make(map[string]string),
			Data:    make([]byte, 1024),
		}
	},
}

func createMessageWithPool(id int) *Message {
	msg := messagePool.Get().(*Message)
	msg.ID = id
	msg.Content = "Hello World"
	// Headers 和 Data 已预分配
	return msg
}

func releaseMessage(msg *Message) {
	// 重置状态
	for k := range msg.Headers {
		delete(msg.Headers, k)
	}
	// Data 保持原容量，只重置长度
	msg.Data = msg.Data[:0]
	messagePool.Put(msg)
}

// 基准测试
func BenchmarkWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg := createMessageWithoutPool(i)
		_ = msg
		// 依赖GC回收
	}
}

func BenchmarkWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg := createMessageWithPool(i)
		_ = msg
		releaseMessage(msg)
	}
}

// func performanceComparison() {
func main() {
	fmt.Println("=== 性能对比 ===")

	const iterations = 100000

	// 无池测试
	start := time.Now()
	for i := 0; i < iterations; i++ {
		msg := createMessageWithoutPool(i)
		_ = msg
	}
	withoutPoolTime := time.Since(start)

	// 有池测试
	start = time.Now()
	for i := 0; i < iterations; i++ {
		msg := createMessageWithPool(i)
		_ = msg
		releaseMessage(msg)
	}
	withPoolTime := time.Since(start)

	fmt.Printf("无池版本: %v\n", withoutPoolTime)
	fmt.Printf("有池版本: %v\n", withPoolTime)
	fmt.Printf("性能提升: %.2fx\n",
		float64(withoutPoolTime)/float64(withPoolTime))
}
