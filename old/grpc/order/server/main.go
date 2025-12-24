package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"testGo/grpc/order"
)

// const addr = ":50052"
var addrList = []string{":50052", ":50053", ":50054"}

func main() {
	Test()
	var wg sync.WaitGroup
	for _, addr := range addrList {
		wg.Add(1)

		go func(val string) {
			defer wg.Done()
			startServer(val)
		}(addr)
	}

	wg.Wait()
}

func Test() {

	jsonstr := `{"id":"1234","items":["aaaaa","bbbb"],"description":"描述","price":123.44,"destination":"订单目的地"}`
	pb := &order.Order{}
	if err := jsonpb.UnmarshalString(jsonstr, pb); err == nil {
		fmt.Errorf("err:%s\n", err)
	}
	fmt.Printf("pb.say:%s\n", pb.Id)
}

// 启动服务
func startServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	// s := grpc.NewServer()

	// 注册拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(orderUnaryServerInterceptor),
		grpc.StreamInterceptor(orderStreamServerInterceptor),
	)

	// 注册订单服务
	orderServer := &OrderServer{orderMap: make(map[string]*order.Order)}
	InitSampleData(orderServer.orderMap)
	order.RegisterOrderManagementServer(s, orderServer)

	// 注册问候服务
	order.RegisterGreeterServiceServer(s, &GreeterServer{})

	log.Println("start gRPC listen on port " + addr)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

// 一元拦截器
func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	// 前置处理
	log.Println("==========[服务端一元拦截器] start ===========", info.FullMethod)

	// 完成方法的正常执行
	res, err = handler(ctx, req)

	// 后置处理
	log.Printf("After method call, res = %+v\n", res)
	log.Println("==========[服务端一元拦截器] end ===========\n\n", info.FullMethod)

	return
}

// WrappedServerStream 服务端流拦截器
type WrappedServerStream struct {
	grpc.ServerStream
}

func (w *WrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("[order stream server interceptor] send a msg : %+v", m)
	return w.ServerStream.SendMsg(m)
}

func (w *WrappedServerStream) RecvMsg(m interface{}) error {
	log.Printf("[order stream server interceptor] recv a msg : %+v", m)
	return w.ServerStream.RecvMsg(m)
}

func NewWrappedServerStream(s grpc.ServerStream) *WrappedServerStream {
	return &WrappedServerStream{s}
}

func orderStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("=========[服务端流拦截器]start %s ========= \n", info.FullMethod)

	// 执行方法
	err := handler(srv, NewWrappedServerStream(ss))
	if err != nil {
		log.Println("handle method err.", err)
	}

	log.Printf("=========[服务端流拦截器]end=========\n\n")
	return nil
}
