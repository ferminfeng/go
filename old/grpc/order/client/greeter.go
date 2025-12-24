package main

import (
	"context"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"testGo/grpc/order"
)

// SayHello ...
func SayHello(ctx context.Context, client order.GreeterServiceClient) {
	log.Println("SayHello")

	_, err := client.SayHello(ctx, &wrappers.StringValue{Value: "roseduan"})
	if err != nil {
		log.Println("call greeter server [say hello] err.", err)
	}
}
