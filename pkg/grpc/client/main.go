package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mikhailshtv/proto_api/pkg/grpc/v1/orders_api"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := orders_api.NewOrderServiceClient(conn)

	createOrder(client)
	getOrder(client)
	getOrders(client)
	editOrder(client)
	deleteOrder(client)
}

func getOrder(client orders_api.OrderServiceClient) {
	order, err := client.GetOrder(context.Background(), &orders_api.OrderActionByIdRequest{Id: 1})
	if err != nil {
		log.Fatalf("error get order request: %s", err.Error())
	}
	orderJson, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Ошибка конвертации в JSON")
	}
	fmt.Println(string(orderJson))
}

func getOrders(client orders_api.OrderServiceClient) {
	orders, err := client.GetOrders(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error get orders request: %s", err.Error())
	}

	ordersJson, err := json.Marshal(orders)
	if err != nil {
		log.Fatal("Ошибка конвертации в JSON")
	}
	fmt.Println(string(ordersJson))
}

func createOrder(client orders_api.OrderServiceClient) {
	product := orders_api.Product{
		Id:            3,
		Code:          140985,
		Quantity:      5553,
		Name:          "Computer",
		PurchasePrice: 5000000,
		SalePrice:     14000000,
	}
	products := []*orders_api.Product{}
	products = append(products, &product)
	order, err := client.CreateOrder(context.Background(), &orders_api.OrderCreateRequest{Products: products})
	if err != nil {
		log.Fatalf("error create order request: %s", err.Error())
	}

	orderJson, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Ошибка конвертации в JSON")
	}
	fmt.Println(string(orderJson))
}

func editOrder(client orders_api.OrderServiceClient) {
	product1 := orders_api.Product{
		Id:            3,
		Code:          140985,
		Quantity:      5553,
		Name:          "Computer",
		PurchasePrice: 5000000,
		SalePrice:     14000000,
	}
	product2 := orders_api.Product{
		Id:            2,
		Code:          137207,
		Quantity:      191,
		Name:          "Bacon",
		PurchasePrice: 34000,
		SalePrice:     137207,
	}
	products := []*orders_api.Product{}
	products = append(products, &product1)
	products = append(products, &product2)
	order, err := client.EditOrder(context.Background(), &orders_api.OrderEditRequest{Id: 1, Products: products})
	if err != nil {
		log.Printf("error edit order request: %s", err.Error())
	}
	orderJson, err := json.Marshal(order)
	if err != nil {
		log.Printf("Ошибка конвертации в JSON")
	}
	fmt.Println(string(orderJson))
}

func deleteOrder(client orders_api.OrderServiceClient) {
	result, err := client.DeleteOrder(context.Background(), &orders_api.OrderActionByIdRequest{Id: 1})
	if err != nil {
		log.Printf("error delete order request: %s", err.Error())
	}
	fmt.Printf("status: %s, message: %s", result.Status, result.Message)
}
