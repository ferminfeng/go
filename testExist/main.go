package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Task ...
type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func main() {
	task := &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Second * 2),
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go task.Run()

	// 主程序在监听`Ctrl+C`消息，收到之后，会把task任务stop
	select {
	case sig := <-c:

		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.Stop()
	}
}

// Run ...
func (t *Task) Run() {
	for {
		select {
		// 这里收到了退出消息之后，就不不会再调用handle了
		case <-t.closed:
			fmt.Println("close .....")
			return
		case <-t.ticker.C:
			t.wg.Add(1)
			fmt.Println("add handle.....")
			go handle(t)
		}
	}
}

// Stop ...
func (t *Task) Stop() {
	fmt.Println("got close sig.....")
	close(t.closed)
	// 在这里会等待所有的协程都退出
	t.wg.Wait()
	fmt.Println("all goroutine　done...")
}

func handle(task *Task) {
	defer task.wg.Done()
	for i := 0; i < 5; i++ {
		fmt.Print("#")
		st := RandInt64(1, 5)
		time.Sleep(time.Second * time.Duration(st))
	}

	fmt.Println()
}

// RandInt64 ...
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
