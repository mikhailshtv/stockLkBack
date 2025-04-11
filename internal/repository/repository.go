package repository

import "golang/stockLkBack/internal/model"

var Orders = []any{}
var Products = []any{}

func CheckAndSaveEntity(entity any) {
	switch entity.(type) {
	case model.Order:
		Orders = append(Orders, entity)
	case model.Product:
		Products = append(Products, entity)
	}
}
