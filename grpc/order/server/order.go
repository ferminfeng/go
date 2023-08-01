package main

import (
	"context"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"strings"
	"testGo/grpc/order"
	"time"
)

type OrderServer struct {
	orderMap map[string]*order.Order
}

// InitSampleData 初始化添加一些订单数据
func InitSampleData(orderMap map[string]*order.Order) {
	orderMap["102"] = &order.Order{Id: "102", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["103"] = &order.Order{Id: "103", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["104"] = &order.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["105"] = &order.Order{Id: "105", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["106"] = &order.Order{Id: "106", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 300.00}
}

// AddOrder 添加订单
func (s *OrderServer) AddOrder(ctx context.Context, req *order.Order) (resp *wrappers.StringValue, err error) {

	// 发送服务端元数据
	md := metadata.New(map[string]string{"location": "San Jose", "timestamp": time.Now().Format(time.StampNano)})
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		log.Println("send header err")
	}

	// 获取客户端元数据
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		log.Println("failed to get metadata")
	} else {
		log.Printf("来自客户端的元数据 : %+v\n", md)
	}

	resp = &wrappers.StringValue{}
	if s.orderMap == nil {
		s.orderMap = make(map[string]*order.Order)
	}

	// time.Sleep(3 * time.Second)

	v4, err := uuid.NewV4()
	if err != nil {
		return resp, status.Errorf(codes.Internal, "gen uuid err", err)
	}
	id := v4.String()
	req.Id = id
	s.orderMap[id] = req
	resp.Value = id
	return
}

// GetOrder 获取订单
func (s *OrderServer) GetOrder(ctx context.Context, req *wrappers.StringValue) (resp *order.Order, err error) {

	resp = &order.Order{}
	id := req.Value
	var exist bool
	if resp, exist = s.orderMap[id]; !exist {
		err = status.Error(codes.NotFound, "order not found id = "+id)
		return
	}
	return
}

// SearchOrder 搜索订单
func (s *OrderServer) SearchOrder(searchKey *wrappers.StringValue, stream order.OrderManagement_SearchOrderServer) (err error) {
	for _, val := range s.orderMap {
		for _, item := range val.Items {
			if strings.Contains(item, searchKey.Value) {
				err = stream.Send(val)
				if err != nil {
					log.Println("stream send order err.", err)
					return
				}
				break
			}
		}
	}
	return
}

// UpdateOrder 更新订单
func (s *OrderServer) UpdateOrder(stream order.OrderManagement_UpdateOrderServer) (err error) {
	updatedIds := "updated order ids : "
	for {
		val, err := stream.Recv()
		if err == io.EOF { // 完成读取订单流
			// 向客户端发送消息
			return stream.SendAndClose(&wrappers.StringValue{Value: updatedIds})
		}
		s.orderMap[val.Id] = val
		log.Println("[server]update the order : ", val.Id)
		updatedIds += val.Id + ", "
	}
}

// ProcessOrder 处理订单
func (s *OrderServer) ProcessOrder(stream order.OrderManagement_ProcessOrderServer) (err error) {
	var combinedShipmentMap = make(map[string]*order.CombinedShipment)
	for {
		val, err := stream.Recv() // 接收从客户端发送来的订单
		if err == io.EOF {        // 接收完毕，返回结果
			for _, shipment := range combinedShipmentMap {
				if err := stream.Send(shipment); err != nil {
					log.Println("[server] process finished!")
					return err
				}
			}
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		if val != nil {
			orderId := val.Value
			log.Printf("[server]reading order : %+v\n", orderId)

			if _, exist := s.orderMap[orderId]; !exist {
				log.Printf("[server]订单不存在 : %+v\n", orderId)
				continue
			}

			dest := s.orderMap[orderId].Destination
			shipment, exist := combinedShipmentMap[dest]
			if exist {
				ord := s.orderMap[orderId]
				shipment.OrderList = append(shipment.OrderList, ord)
				combinedShipmentMap[dest] = shipment
			} else {
				comShip := &order.CombinedShipment{Id: "cmb - " + (s.orderMap[orderId].Destination), Status: "Processed!"}
				ord := s.orderMap[orderId]
				comShip.OrderList = append(comShip.OrderList, ord)
				combinedShipmentMap[dest] = comShip
				log.Println(len(comShip.OrderList), comShip.GetId())
			}
		}
	}
	return
}
