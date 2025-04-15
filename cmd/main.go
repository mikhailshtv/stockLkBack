package main

import (
	"encoding/json"
	"fmt"
	"golang/stockLkBack/internal/repository"
	"golang/stockLkBack/internal/service"
	"time"
)

func main() {
	go inerval()
	time.Sleep(time.Second * 20)
}

func inerval() {
	for range time.Tick(time.Second * 1) {
		entity := service.NewEntity()
		repository.CheckAndSaveEntity(entity)
		entityJSON, _ := json.Marshal(entity)
		ordersJSON, _ := json.Marshal(repository.OrdersStruct.Entities)
		productsJSON, _ := json.Marshal(repository.ProductsStruct.Entities)
		usersJSON, _ := json.Marshal(repository.UsersStruct.Entities)
		fmt.Printf("Entity: %v\n", string(entityJSON))
		fmt.Printf("Orders: %v\n", string(ordersJSON))
		fmt.Printf("Products: %v\n", string(productsJSON))
		fmt.Printf("Users: %v\n", string(usersJSON))
	}
}
