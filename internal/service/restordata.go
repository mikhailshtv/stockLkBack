package service

import (
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/repository"
	"log"
	"os"
)

func RestoreData() {
	if _, err := os.Stat("./assets/orders.json"); !errors.Is(err, os.ErrNotExist) {
		repository.OrdersStruct.RestoreFromFile("./assets/orders.json")
		ordersJSON, err := json.Marshal(repository.OrdersStruct.Entities)
		if err != nil {
			log.Fatalf("Ошибка сериализации %v", err.Error())
		}
		log.Printf("восстановленные данные orders: %v\n", string(ordersJSON))
	}
	if _, err := os.Stat("./assets/products.json"); !errors.Is(err, os.ErrNotExist) {
		repository.ProductsStruct.RestoreFromFile("./assets/products.json")
		productsJSON, err := json.Marshal(repository.ProductsStruct.Entities)
		if err != nil {
			log.Fatalf("Ошибка сериализации %v", err.Error())
		}
		log.Printf("восстановленные данные products: %v\n", string(productsJSON))
	}
	if _, err := os.Stat("./assets/users.json"); !errors.Is(err, os.ErrNotExist) {
		repository.UsersStruct.RestoreFromFile("./assets/users.json")
		usersJSON, err := json.Marshal(repository.OrdersStruct.Entities)
		if err != nil {
			log.Fatalf("Ошибка сериализации %v", err.Error())
		}
		log.Printf("восстановленные данные users: %v\n", string(usersJSON))
	}
}
