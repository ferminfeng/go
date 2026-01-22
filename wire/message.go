package main

import (
	"errors"
	"fmt"
	"time"
)

type Message struct {
	Content string
	Code    int
}

// 接收参数作为消息内容
//func NewMessage(phrase string) Message {
//	return Message{
//		Content: phrase,
//	}
//}

type Greeter struct {
	Message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

type Event struct {
	Greeter Greeter
}

// 增加返回错误信息
func NewEvent(g Greeter) (Event, error) {
	// 模拟创建 Event 报错
	if time.Now().Unix()%2 == 0 {
		return Event{}, errors.New("new event error")
	}
	return Event{Greeter: g}, nil
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}
