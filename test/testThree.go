package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		second := time.Now().Second()
		fmt.Println("second:", second)
		time.Sleep(time.Second * 1)

		if second%5 == 0 {
			fmt.Println("continue:")
			continue
		}

		fmt.Println("second: 下一步")
		if second == 1 {
			fmt.Println("break:")
			break
		}

	}
}
