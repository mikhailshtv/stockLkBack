package repository

import (
	"golang/stockLkBack/internal/model"
)

type Entity[T model.Order | model.Product | model.User] struct {
	Entities []*T
}

func (entity *Entity[T]) AppendEntity(v T) {
	entity.Entities = append(entity.Entities, &v)
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
