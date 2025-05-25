package handler

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(ctx *gin.Context) {
	var order model.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	repository.CheckAndSaveEntity(order)
	ctx.JSON(http.StatusOK, order)
}
