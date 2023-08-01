package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"testGo/grpc/order"
)

const address = "localhost:50052"

func main() {
	// conn, err := grpc.Dial(address, grpc.WithInsecure())

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		// 注册拦截器
		grpc.WithUnaryInterceptor(UnaryClientOrderInterceptor),
		grpc.WithStreamInterceptor(StreamClientOrderInterceptor),
	)

	if err != nil {
		log.Println("did not connect.", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	ctx := context.Background()
	//
	// // 使用带有截止时间的context
	// ctx, cancel := context.WithDeadline(
	// 	context.Background(),
	// 	// 适当调整截止时间观察不同的调用效果
	// 	time.Now().Add(2*time.Second))
	//
	// defer cancel()

	client := order.NewOrderManagementClient(conn)

	// 问候服务客户端
	greeterClient := order.NewGreeterServiceClient(conn)

	fmt.Println("-----------unary rpc-------------")

	// 取消RPC请求
	// cancelRpcRequest(client)

	// 添加订单
	id := AddOrder(ctx, client)

	// 获取订单
	GetOrder(ctx, client, id)

	// // 搜索订单
	// // SearchOrder(ctx, client)
	//
	// // 更新订单
	// // UpdateOrder(ctx, client)
	//
	// // 处理订单
	// ProcessOrder(ctx, client)

	// 问候服务
	SayHello(ctx, greeterClient)
}

// UnaryClientOrderInterceptor 客户端一元拦截器
func UnaryClientOrderInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	log.Println("=========[客户端一元拦截器]start ========= ", method)
	err = invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Println("invoke method err.", err)
	}
	log.Println("=========[client interceptor] end. reply : ", reply)
	log.Println("=========[客户端一元拦截器]end =========\n\n", method)
	return
}

// WrappedClientStream 客户端流拦截器
type WrappedClientStream struct {
	grpc.ClientStream
}

func (w *WrappedClientStream) SendMsg(m interface{}) error {
	log.Printf("===========[client interceptor] send msg : %+v", m)
	return w.ClientStream.SendMsg(m)
}

func (w *WrappedClientStream) RecvMsg(m interface{}) error {
	log.Printf("============[client interceptor] recv msg : %+v", m)
	return w.ClientStream.RecvMsg(m)
}

func NewWrappedClientStream(s grpc.ClientStream) *WrappedClientStream {
	return &WrappedClientStream{s}
}

func StreamClientOrderInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("===========[客户端流拦截器]start, method = %+v\n", method)
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return NewWrappedClientStream(clientStream), nil
}
