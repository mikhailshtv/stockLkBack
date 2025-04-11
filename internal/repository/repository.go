package repository

import (
	"golang/stockLkBack/internal/model"
	"reflect"
)

var Orders = []model.Order{}
var Products = []model.Product{}
var Users = []model.User{}

func CheckAndSaveEntity(entity any) {
	switch entity.(type) {
	case model.Order:
		Orders = append(Orders, reflect.ValueOf(entity).Interface().(model.Order))
	case model.Product:
		Products = append(Products, reflect.ValueOf(entity).Interface().(model.Product))
	case model.User:
		Users = append(Users, reflect.ValueOf(entity).Interface().(model.User))
	}
}
