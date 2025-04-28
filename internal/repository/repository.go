package repository

import (
	"encoding/json"
	"golang/stockLkBack/internal/model"
	"log"
	"sync"
	"time"
)

type Entity[T model.Order | model.Product | model.User] struct {
	Mu       sync.RWMutex
	Entities []*T
}

func (entity *Entity[T]) AppendEntity(v T) {
	entity.Mu.Lock()
	entity.Entities = append(entity.Entities, &v)
	entity.Mu.Unlock()
}

var OrdersStruct = Entity[model.Order]{}
var ProductsStruct = Entity[model.Product]{}
var UsersStruct = Entity[model.User]{}

func CheckAndSaveEntity(entity any) {
	switch v := entity.(type) {
	case model.Order:
		OrdersStruct.AppendEntity(v)
	case model.Product:
		ProductsStruct.AppendEntity(v)
	case model.User:
		UsersStruct.AppendEntity(v)
	}
}

func LogAddedEntities() {
	var ordersLength int
	var productsLength int
	var userLength int
	for range time.Tick(time.Millisecond * 200) {
		if len(OrdersStruct.Entities) != ordersLength {
			OrdersStruct.Mu.RLock()
			ordersJSON, err := json.Marshal(OrdersStruct.Entities[ordersLength:len(OrdersStruct.Entities)])
			OrdersStruct.Mu.RUnlock()
			if err == nil {
				log.Printf("Orders: %v\n", string(ordersJSON))
			}
		}
		if len(ProductsStruct.Entities) != productsLength {
			ProductsStruct.Mu.RLock()
			productsJSON, err := json.Marshal(ProductsStruct.Entities[productsLength:len(ProductsStruct.Entities)])
			ProductsStruct.Mu.RUnlock()
			if err == nil {
				log.Printf("Products: %v\n", string(productsJSON))
			}
		}
		if len(UsersStruct.Entities) != userLength {
			UsersStruct.Mu.RLock()
			usersJSON, err := json.Marshal(UsersStruct.Entities[userLength:len(UsersStruct.Entities)])
			UsersStruct.Mu.RUnlock()
			if err == nil {
				log.Printf("Users: %v\n", string(usersJSON))
			}
		}
		ordersLength = len(OrdersStruct.Entities)
		productsLength = len(ProductsStruct.Entities)
		userLength = len(UsersStruct.Entities)
	}
}
