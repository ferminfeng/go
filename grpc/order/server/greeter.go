package main

import (
	"context"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"log"
)

type GreeterServer struct {
}

func (s *GreeterServer) SayHello(ctx context.Context, req *wrappers.StringValue) (resp *wrappers.StringValue, err error) {
	resp = &wrappers.StringValue{}
	log.Printf("[greeter server]Hello, %s, this is greeter server.", req.Value)
	return
}
