package main

import (
	"fmt"
	nats "github.com/nats-io/nats.go"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	sub()
}

func sub() {
	// 连接Nats服务器
	nc, _ := nats.Connect("nats://127.0.0.1:4222")

	// 发布-订阅 模式，异步订阅 test1
	_, _ = nc.Subscribe("test1", func(m *nats.Msg) {
		fmt.Printf("test1 Received a message: %s\n", string(m.Data))
	})

	// 队列 模式，订阅 test2， 队列为queue, test2 发向所有队列，同一队列只有一个能收到消息
	_, _ = nc.QueueSubscribe("test2", "queue", func(msg *nats.Msg) {
		fmt.Printf("test2 Queue a message: %s\n", string(msg.Data))
	})

	// 请求-响应， 响应 test3 消息。
	_, _ = nc.Subscribe("test3", func(m *nats.Msg) {
		fmt.Printf("test3 Reply a message: %s\n", string(m.Data))
		_ = nc.Publish(m.Reply, []byte("I can help for test3!!"))
	})

	// 持续发送不需要关闭
	//_ = nc.Drain()

	// 关闭连接
	//nc.Close()

	// 阻止进程结束而收不到消息
	signal.Ignore(syscall.SIGHUP)
	runtime.Goexit()
}
