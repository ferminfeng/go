package main

import (
	"fmt"
	"os"
)

func main() {
	target := "World"
	if len(os.Args) > 1 { /* os.Args是一个参数切片 */
		target = os.Args[1]
	}
	fmt.Println("Hello", target)
}
