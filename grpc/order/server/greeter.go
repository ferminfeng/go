package main

import (
	"context"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"log"
)

type GreeterServer struct {
}

// SayHello ...
func (s *GreeterServer) SayHello(ctx context.Context, req *wrappers.StringValue) (resp *wrappers.StringValue, err error) {
	resp = &wrappers.StringValue{}
	log.Printf("[问候服务]Hello, %s, this is greeter server.", req.Value)
	return
}
