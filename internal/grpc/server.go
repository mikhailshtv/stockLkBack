package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/handler"
	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"

	"github.com/mikhailshtv/proto_api/pkg/grpc/v1/orders_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	orders_api.UnimplementedOrderServiceServer
	handler *handler.Handler
}

func (s *server) GetOrders(
	_ context.Context,
	req *orders_api.OrderGetAllRequest,
) (*orders_api.GetOrdersResponse, error) {
	ordersAll, err := s.handler.Services.Order.GetAll(req.UserId, model.UserRole(req.Role))
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка получения списка заказов")
		if err != nil {
			log.Println(err.Error())
		}
	}
	ordersJSON, err := json.Marshal(ordersAll)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	var orders []*orders_api.Order
	err = json.Unmarshal(ordersJSON, &orders)
	if err != nil {
		log.Println(err.Error())
	}
	return &orders_api.GetOrdersResponse{
		Orders: orders,
	}, nil
}

func (s *server) GetOrder(
	_ context.Context,
	req *orders_api.OrderGetByIdRequest,
) (*orders_api.GetOrderResponse, error) {
	orderID := req.GetId()
	if orderID <= 0 {
		err := status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}

	receivedOrder, err := s.handler.Services.Order.GetByID(orderID, req.UserId, model.UserRole(req.Role))
	if err != nil {
		if err.Error() == repository.NotFoundErrorMessage {
			err = status.Errorf(codes.NotFound, "Объект не найден")
			if err != nil {
				log.Println(err.Error())
			}
			return nil, err
		}
		err = status.Errorf(codes.Internal, "%s", err.Error())
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}
	orderJSON, err := json.Marshal(receivedOrder)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	var order *orders_api.Order
	err = json.Unmarshal(orderJSON, &order)
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
	userID := req.UserId
	productsJSON, err := json.Marshal(products)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	var orderReq model.OrderRequestBody

	err = json.Unmarshal(productsJSON, &orderReq.Products)
	if err != nil {
		log.Println(err.Error())
	}
	order, err := s.handler.Services.Order.Create(orderReq, userID)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при создании заказа")
		if err != nil {
			log.Println(err.Error())
		}
	}

	protoProductsJSON, err := json.Marshal(order.Products)
	if err != nil {
		log.Println(err.Error())
	}
	var protoProducts []*orders_api.Product
	err = json.Unmarshal(protoProductsJSON, &protoProducts)
	if err != nil {
		log.Println(err.Error())
	}
	return &orders_api.Order{
		Id:               order.ID,
		Number:           order.Number,
		TotalCost:        order.TotalCost,
		CreatedDate:      timestamppb.New(order.CreatedDate),
		LastModifiedDate: timestamppb.New(order.LastModifiedDate),
		Status:           &orders_api.OrderStatus{Key: order.Status.Key, DisplayName: order.Status.DisplayName},
		Products:         protoProducts,
	}, nil
}

func (s *server) EditOrder(
	_ context.Context,
	req *orders_api.OrderEditRequest,
) (*orders_api.Order, error) {
	orderID := req.GetId()
	userID := req.UserId
	if orderID <= 0 {
		err := status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}

	var orderReq model.OrderRequestBody
	orderReqJSON, err := json.Marshal(req)
	if err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка при конвертации в JSON")
		if err != nil {
			log.Println(err.Error())
		}
	}
	if err := json.Unmarshal(orderReqJSON, &orderReq); err != nil {
		log.Println(err.Error())
		err = status.Errorf(codes.Internal, "Ошибка десериализации")
		if err != nil {
			log.Println(err.Error())
		}
	}

	order, err := s.handler.Services.Order.Update(orderID, orderReq, userID)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == repository.NotFoundErrorMessage {
			err = status.Errorf(codes.NotFound, "Объект не найден")
			if err != nil {
				log.Println(err.Error())
			}
			return nil, errors.New("объект не найден")
		}
		err = status.Errorf(codes.Internal, "%s", err.Error())
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}

	protoProductsJSON, err := json.Marshal(order.Products)
	if err != nil {
		log.Println(err.Error())
	}
	var protoProducts []*orders_api.Product
	err = json.Unmarshal(protoProductsJSON, &protoProducts)
	if err != nil {
		log.Println(err.Error())
	}
	return &orders_api.Order{
		Id:               order.ID,
		Number:           order.Number,
		TotalCost:        order.TotalCost,
		CreatedDate:      timestamppb.New(order.CreatedDate),
		LastModifiedDate: timestamppb.New(order.LastModifiedDate),
		Status:           &orders_api.OrderStatus{Key: order.Status.Key, DisplayName: order.Status.DisplayName},
		Products:         protoProducts,
	}, nil
}

func (s *server) DeleteOrder(
	_ context.Context,
	req *orders_api.OrderDeleteRequest,
) (*orders_api.Success, error) {
	orderID := req.GetId()
	userID := req.UserId
	if orderID <= 0 {
		err := status.Errorf(codes.InvalidArgument, "id должен быть больше чем 0")
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
	}

	err := s.handler.Services.Order.Delete(orderID, userID)
	if err != nil {
		if err.Error() == repository.NotFoundErrorMessage {
			err = status.Errorf(codes.NotFound, "Объект не найден")
			if err != nil {
				log.Println(err.Error())
			}
			return nil, errors.New("объект не найден")
		}
		err = status.Errorf(codes.Internal, "%s", err.Error())
		if err != nil {
			log.Println(err.Error())
		}
		return nil, err
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
