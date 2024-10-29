package main

import (
	"fmt"
)

func main() {
	for n := 2; n <= 100; n++ {
		if isPrime(n) == true {
			fmt.Printf("%d 是质数\n", n)
		}
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
