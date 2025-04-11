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
		ordersJSON, _ := json.Marshal(repository.Orders)
		productsJSON, _ := json.Marshal(repository.Products)
		fmt.Printf("Entity: %v\n", string(entityJSON))
		fmt.Printf("Orders: %v\n", string(ordersJSON))
		fmt.Printf("Products: %v\n", string(productsJSON))
	}
}
