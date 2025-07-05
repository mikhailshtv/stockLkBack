package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/handler"
	"golang/stockLkBack/internal/model"
	"log"
	"net"
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
	handler *handler.Handler
}

func (s *server) GetOrders(
	_ context.Context,
	_ *emptypb.Empty,
) (*orders_api.GetOrdersResponse, error) {
	ordersAll, err := s.handler.Services.Order.GetAll()
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка получения списка заказов")
		if err != nil {
			log.Println(err.Error())
		}
	}
	ordersJson, err := json.Marshal(ordersAll)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	var orders []*orders_api.Order
	err = json.Unmarshal(ordersJson, &orders)
	if err != nil {
		log.Println(err.Error())
	}
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
		err := status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}

	receivedOrder, err := s.handler.Services.Order.GetById(int(orderId))
	if err != nil {
		log.Fatal(err.Error())
		if err.Error() == "элемент не найден" {
			err = status.Errorf(codes.NotFound, "Объект не найден")
			if err != nil {
				log.Println(err.Error())
			}
			return nil, err
		} else {
			err = status.Errorf(codes.Internal, err.Error())
			if err != nil {
				log.Println(err.Error())
			}
			return nil, err
		}

	}
	orderJson, err := json.Marshal(receivedOrder)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	var order *orders_api.Order
	err = json.Unmarshal(orderJson, &order)
	if err != nil {
		log.Println(err.Error())
	}
	return &orders_api.GetOrderResponse{
		Order: order,
	}, nil
}

func (s *server) CreateOrder(
	_ context.Context,
	req *orders_api.OrderCreateRequest,
) (*orders_api.Order, error) {
	products := req.Products
	productsJson, err := json.Marshal(products)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	var orderReq model.OrderRequestBody

	err = json.Unmarshal(productsJson, &orderReq.Products)
	if err != nil {
		log.Println(err.Error())
	}
	order, err := s.handler.Services.Order.Create(orderReq)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при создании заказа")
		if err != nil {
			log.Println(err.Error())
		}
	}
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
		err := status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}
	products := req.Products

	var orderReq model.OrderRequestBody
	orderReqJson, err := json.Marshal(req)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	if err := json.Unmarshal(orderReqJson, &orderReq); err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка десериализации")
		if err != nil {
			log.Println(err.Error())
		}
	}

	order, err := s.handler.Services.Order.Update(int(orderId), orderReq)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "элемент не найден" {
			err = status.Errorf(codes.NotFound, "Объект не найден")
			if err != nil {
				log.Println(err.Error())
			}
			return nil, errors.New("объект не найден")
		} else {
			err = status.Errorf(codes.Internal, err.Error())
			if err != nil {
				log.Println(err.Error())
			}
			return nil, err
		}
	}

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

func (s *server) DeleteOrder(
	_ context.Context,
	req *orders_api.OrderActionByIdRequest,
) (*orders_api.Success, error) {
	orderId := req.GetId()
	if orderId <= 0 {
		err := status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}

	err := s.handler.Services.Order.Delete(int(orderId))
	if err != nil {
		if err.Error() == "элемент не найден" {
			err = status.Errorf(codes.NotFound, "Объект не найден")
			if err != nil {
				log.Println(err.Error())
			}
			return nil, errors.New("объект не найден")
		} else {
			err = status.Errorf(codes.Internal, err.Error())
			if err != nil {
				log.Println(err.Error())
			}
			return nil, err
		}
	}

	return &orders_api.Success{
		Status:  "Success",
		Message: "Объект успешно удален",
	}, nil
}

func StartServer(handler *handler.Handler) {
	lis, err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingInterceptor,
		),
	)

	orders_api.RegisterOrderServiceServer(s, &server{handler: handler})
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
