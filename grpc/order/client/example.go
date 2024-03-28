// @Author fermin 2023/8/1 13:42:00
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"testGo/grpc/order"
	"time"
)

func makeRPCs(cc *grpc.ClientConn, n int) {
	client := order.NewEchoServiceClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(client, "test for load balance")
	}
}

func callUnaryEcho(c order.EchoServiceClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &wrappers.StringValue{Value: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Value)
}

type exampleResolverBuilder struct{}

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addressList,
		},
	}

	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

func (r *exampleResolver) start() {
	// addrStrs := r.addrsStore[r.target.Endpoint]
	addrStrs := r.addrsStore[exampleServiceName]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
