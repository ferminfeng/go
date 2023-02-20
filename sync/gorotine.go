package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan int, 1)

	go func() {
		fmt.Println("start gorotine")
		time.Sleep(5 * time.Second)
		done <- 1
		fmt.Println("end gorotine")
	}()

	fmt.Println("waiting")
	<-done
	fmt.Println("done")
}
