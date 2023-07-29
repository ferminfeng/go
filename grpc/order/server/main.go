package main

import (
	"context"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"net"
	"strings"
	"testGo/grpc/order"
)

const port = ":50052"

type server struct {
	orderMap map[string]*order.Order
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	order.RegisterOrderManagementServer(s, &server{})
	log.Println("start gRPC listen on port " + port)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

// AddOrder 添加订单
func (s *server) AddOrder(ctx context.Context, req *order.Order) (resp *wrappers.StringValue, err error) {
	resp = &wrappers.StringValue{}
	if s.orderMap == nil {
		s.orderMap = make(map[string]*order.Order)
	}

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
func (s *server) GetOrder(ctx context.Context, req *wrappers.StringValue) (resp *order.Order, err error) {
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
func (s *server) SearchOrder(searchKey *wrappers.StringValue, stream order.OrderManagement_SearchOrderServer) (err error) {
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
func (s *server) UpdateOrder(stream order.OrderManagement_UpdateOrderServer) (err error) {
	updatedIds := "updated order ids : "
	for {
		val, err := stream.Recv()
		if err == io.EOF { //完成读取订单流
			//向客户端发送消息
			return stream.SendAndClose(&wrappers.StringValue{Value: updatedIds})
		}
		s.orderMap[val.Id] = val
		log.Println("[server]update the order : ", val.Id)
		updatedIds += val.Id + ", "
	}
}

// ProcessOrder 处理订单
func (s *server) ProcessOrder(stream order.OrderManagement_ProcessOrderServer) (err error) {
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
