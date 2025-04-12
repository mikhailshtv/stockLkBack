package repository

import (
	"golang/stockLkBack/internal/model"
)

var Orders = []*model.Order{}
var Products = []*model.Product{}
var Users = []*model.User{}

func CheckAndSaveEntity(entity any) {
	switch v := entity.(type) {
	case model.Order:
		Orders = append(Orders, &v)
	case model.Product:
		Products = append(Products, &v)
	case model.User:
		Users = append(Users, &v)
	}
}
