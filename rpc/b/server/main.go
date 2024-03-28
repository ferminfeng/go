package main

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	log.Println("Arith => Multiply", args)
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	log.Println("Arith => Divide", args)

	if args.B == 0 {
		return errors.New("divide by 0")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	log.Println("开始启动rpc 服务端")
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("serve error:", err)
	}
}
