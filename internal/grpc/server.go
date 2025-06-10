package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"golang/stockLkBack/internal/service"
	"log"
	"net"
	"slices"
	"time"

	"github.com/mikhailshtv/proto_api/pkg/grpc/v1/orders_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	orders_api.UnimplementedOrderServiceServer
}

func (s *server) GetOrders(
	_ context.Context,
	_ *emptypb.Empty,
) (*orders_api.GetOrdersResponse, error) {
	ordersJson, err := json.Marshal(repository.OrdersStruct.Entities)
	if err != nil {
		log.Fatal(err.Error())
		status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
	}
	var orders []*orders_api.Order
	json.Unmarshal(ordersJson, &orders)
	return &orders_api.GetOrdersResponse{
		Orders: orders,
	}, nil
}

func (s *server) GetOrder(
	_ context.Context,
	req *orders_api.OrderActionByIdRequest,
) (*orders_api.GetOrderResponse, error) {
	orderId := req.GetId()
	if orderId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
	}

	for _, v := range repository.OrdersStruct.Entities {
		if v.Id == int(orderId) {
			orderJson, err := json.Marshal(v)
			if err != nil {
				log.Fatal(err.Error())
				status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
			}
			var order *orders_api.Order
			json.Unmarshal(orderJson, &order)
			return &orders_api.GetOrderResponse{
				Order: order,
			}, nil
		}
	}
	status.Errorf(codes.NotFound, "Объект не найден")
	return nil, errors.New("объект не найден")
}

func (s *server) CreateOrder(
	_ context.Context,
	req *orders_api.OrderCreateRequest,
) (*orders_api.Order, error) {
	products := req.Products
	productsJson, err := json.Marshal(products)
	if err != nil {
		log.Fatal(err.Error())
		status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
	}
	var order model.Order

	json.Unmarshal(productsJson, &order.Products)
	service.SetCommonOrderDataOnCreate(&order)
	repository.CheckAndSaveEntity(order)
	return &orders_api.Order{
		Id:               int32(order.Id),
		Number:           int32(order.Number),
		TotalCost:        int32(order.TotalCost),
		CreatedDate:      timestamppb.New(order.CreatedDate),
		LastModifiedDate: timestamppb.New(order.LastModifiedDate),
		Status:           int32(order.Status),
		Products:         products,
	}, nil
}

func (s *server) EditOrder(
	_ context.Context,
	req *orders_api.OrderEditRequest,
) (*orders_api.Order, error) {
	orderId := req.GetId()
	if orderId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
	}
	products := req.Products

	for _, v := range repository.OrdersStruct.Entities {
		if v.Id == int(orderId) {
			productsJson, err := json.Marshal(products)
			if err != nil {
				log.Fatal(err.Error())
				status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
			}
			json.Unmarshal(productsJson, &v.Products)
			v.LastModifiedDate = time.Now().UTC()
			v.TotalCost = 0
			for _, product := range v.Products {
				v.TotalCost += product.SalePrice
			}
			repository.OrdersStruct.SaveToFile("./assets/orders.json")

			return &orders_api.Order{
				Id:               int32(v.Id),
				Number:           int32(v.Number),
				TotalCost:        int32(v.TotalCost),
				CreatedDate:      timestamppb.New(v.CreatedDate),
				LastModifiedDate: timestamppb.New(v.LastModifiedDate),
				Status:           int32(v.Status),
				Products:         products,
			}, nil
		}
	}
	status.Errorf(codes.NotFound, "Объект не найден")
	return nil, errors.New("объект не найден")
}

func (s *server) DeleteOrder(
	_ context.Context,
	req *orders_api.OrderActionByIdRequest,
) (*orders_api.Success, error) {
	orderId := req.GetId()
	if orderId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
	}

	for i, v := range repository.OrdersStruct.Entities {
		if v.Id == int(orderId) {
			repository.OrdersStruct.Entities = slices.Delete(repository.OrdersStruct.Entities, i, i+1)
			repository.OrdersStruct.EntitiesLen = len(repository.OrdersStruct.Entities)
			repository.OrdersStruct.SaveToFile("./assets/orders.json")

			return &orders_api.Success{
				Status:  "Success",
				Message: "Объект успешно удален",
			}, nil
		}
	}
	status.Errorf(codes.NotFound, "Объект не найден")
	return nil, errors.New("объект не найден")
}

func StartServer() {
	lis, err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingInterceptor,
		),
	)

	orders_api.RegisterOrderServiceServer(s, &server{})
	reflection.Register(s)

	log.Println("Server is running at :5001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	st, _ := status.FromError(err)

	var reqJSON, respJSON string

	if m, ok := req.(proto.Message); ok {
		b, _ := protojson.Marshal(m)
		reqJSON = string(b)
	} else {
		reqJSON = "<non-proto request>"
	}

	if m, ok := resp.(proto.Message); ok && resp != nil {
		b, _ := protojson.Marshal(m)
		respJSON = string(b)
	} else {
		respJSON = "<non-proto response or nil>"
	}

	log.Printf(
		"[gRPC] method=%s status=%s error=%v duration=%s request=%s response=%s",
		info.FullMethod,
		st.Code(),
		err,
		time.Since(start),
		reqJSON,
		respJSON,
	)

	return resp, err
}
