package main

import (
	"log"
	"net"
	"net/rpc"
)

func main() {
	port := "9999"
	rpc.RegisterName("HelloService", new(HelloService))
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	rpc.Accept(listener)

}

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	log.Println("got req", request)
	return nil
}
