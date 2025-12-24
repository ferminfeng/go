package main

import (
	"context"
	"fmt"
	"time"
)

func processRequest(ctx context.Context, id int) {
	select {
	case <-time.After(2 * time.Second): // 模拟耗时操作
		fmt.Printf("Request %d completed\n", id)
	case <-ctx.Done(): // 监听取消信号
		fmt.Printf("Request %d cancelled: %v\n", id, ctx.Err())
		return
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go processRequest(ctx, 1)
	time.Sleep(3 * time.Second) // 等待goroutine完成
}
