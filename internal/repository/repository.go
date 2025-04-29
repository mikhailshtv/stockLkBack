package repository

import (
	"golang/stockLkBack/internal/model"
	"sync"
)

type Entity[T model.Order | model.Product | model.User] struct {
	Mu          sync.RWMutex
	Entities    []*T
	EntitiesLen int
}

func (entity *Entity[T]) AppendEntity(v T) {
	entity.Mu.Lock()
	entity.Entities = append(entity.Entities, &v)
	entity.Mu.Unlock()
}

func (entity *Entity[T]) SavedEntities() []*T {
	entity.Mu.RLock()
	defer entity.Mu.RUnlock()
	return entity.Entities[entity.EntitiesLen:]
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
