package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {

	client, err := rpc.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Call("HelloService.Hello", "哈哈哈哈", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
