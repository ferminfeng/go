package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"testGo/grpc/order"
	"time"
)

// AddOrder 添加订单
func AddOrder(ctx context.Context, client order.OrderManagementClient) string {
	log.Println("添加订单")

	odr := &order.Order{
		Description: "a new order for test-1",
		Price:       12322.232,
		Destination: "Shanghai",
		Items:       []string{"doll", "22", "33", "Apple"},
	}
	// 接收从服务端发送过来的metadata信息
	var header metadata.MD

	val, err := client.AddOrder(ctx, odr, grpc.Header(&header))
	if err != nil {
		log.Println("add order fail.", err)
		return ""
	}
	log.Println("add order success.id = ", val.String())

	log.Printf("接受服务端元数据 : %+v\n", header)

	fmt.Println("")
	return val.Value
}

// GetOrder 获取订单
func GetOrder(ctx context.Context, client order.OrderManagementClient, id string) {
	log.Println("获取订单")

	val, err := client.GetOrder(ctx, &wrappers.StringValue{Value: id})
	if err != nil {
		log.Println("get order err.", err)
		return
	}

	log.Printf("get order succes. order = %+v", val)

	fmt.Println("")
}

// SearchOrder 搜索订单
func SearchOrder(ctx context.Context, client order.OrderManagementClient) {
	log.Println("搜索订单")

	searchKey := "Apple"
	searchStream, _ := client.SearchOrder(ctx, &wrappers.StringValue{Value: searchKey})
	for {
		val, err := searchStream.Recv()
		if err == io.EOF { // 服务端没有数据了
			break
		}
		log.Printf("search order from server : %+v", val)
	}

	fmt.Println("")
	return
}

// UpdateOrder 更新订单
func UpdateOrder(ctx context.Context, client order.OrderManagementClient) {
	log.Println("更新订单")

	updateStream, _ := client.UpdateOrder(ctx)
	order1 := &order.Order{Id: "103", Items: []string{"Apple Watch S6"}, Destination: "San Jose, CA", Price: 4400.00}
	order2 := &order.Order{Id: "105", Items: []string{"Amazon Kindle"}, Destination: "San Jose, CA", Price: 330.00}

	// 更新订单1
	if err := updateStream.Send(order1); err != nil {
		log.Println("send order err.", err)
	}

	// 更新订单2
	if err := updateStream.Send(order2); err != nil {
		log.Println("send order err.", err)
	}

	// 关闭流并接收响应
	recv, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Println("close and recv err.", err)
		return
	}
	log.Printf("the update result : %+v", recv)

	fmt.Println("")
}

// ProcessOrder 处理订单
func ProcessOrder(ctx context.Context, client order.OrderManagementClient) {
	log.Println("处理订单")

	processStream, _ := client.ProcessOrder(ctx)

	// 发送两个订单处理
	if err := processStream.Send(&wrappers.StringValue{Value: "103"}); err != nil {
		log.Println("send order err.", err)
	}

	if err := processStream.Send(&wrappers.StringValue{Value: "105"}); err != nil {
		log.Println("send order err.", err)
	}

	chn := make(chan struct{})
	// 异步接收服务端的结果
	go processResultFromServer(processStream, chn)

	// 再发送一个订单
	if err := processStream.Send(&wrappers.StringValue{Value: "106"}); err != nil {
		log.Println("send order err.", err)
	}
	// 发送完毕后记得关闭
	if err := processStream.CloseSend(); err != nil {
		log.Println("close send err.", err)
	}

	<-chn

	fmt.Println("")
}

// 从服务端获取处理的结果
func processResultFromServer(stream order.OrderManagement_ProcessOrderClient, chn chan struct{}) {
	defer close(chn)
	for {
		shipment, err := stream.Recv()
		if err == io.EOF {
			log.Println("[client]结束从服务端接收数据")
			break
		}
		log.Printf("[client]server process result : %+v\n", shipment)
	}
}

// 取消RPC请求
func cancelRpcRequest(client order.OrderManagementClient) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	done := make(chan string)
	go func() {
		var id string
		defer func() {
			fmt.Println("结束执行, id = ", id)
			done <- id
		}()

		time.Sleep(2 * time.Second)
		id = AddOrder(ctx, client)
		log.Println("添加订单成功, id = ", id)
	}()

	// 等待一秒后取消
	time.Sleep(time.Second)
	cancelFunc()

	<-done
}
