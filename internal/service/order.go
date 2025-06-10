package service

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"time"
)

func SetCommonOrderDataOnCreate(order *model.Order) {
	order.Id = repository.OrdersStruct.EntitiesLen + 1

	if repository.OrdersStruct.EntitiesLen > 0 {
		lastOrder := repository.OrdersStruct.Entities[repository.OrdersStruct.EntitiesLen-1]
		order.Number = lastOrder.Number + 1
	} else {
		order.Number = 1
	}
	order.CreatedDate = time.Now().UTC()
	order.LastModifiedDate = time.Now().UTC()
	order.Status = model.Active
	totalCost := 0
	for _, product := range order.Products {
		totalCost += product.SalePrice
	}
	order.TotalCost = totalCost
}
