package main

import "fmt"

func main() {
	Hello("asasa")
}

func Hello(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
