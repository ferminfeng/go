package main

import (
	"fmt"
	"sync"
)

// 三个协程，顺序输出1/2/3

func main() {
	var mg sync.WaitGroup
	mg.Add(3)

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		defer mg.Done()
		fmt.Println(1)

		ch1 <- 1
	}()

	go func() {
		defer mg.Done()
		<-ch1
		fmt.Println(2)

		ch2 <- 1
	}()

	go func() {
		defer mg.Done()
		<-ch2
		fmt.Println(3)

		ch3 <- 1
	}()

	<-ch3
	mg.Wait()

}
