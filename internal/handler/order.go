package handler

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"log"
	"net/http"
	"slices"
	"strconv"
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

func EditOrder(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.OrdersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			var order model.Order
			if err := ctx.ShouldBindJSON(&order); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			v.Products = order.Products
			v.LastModifiedDate = time.Now().UTC()
			v.TotalCost = 0
			for _, product := range order.Products {
				v.TotalCost += product.SalePrice
			}
			repository.OrdersStruct.SaveToFile("./assets/orders.json")
			ctx.JSON(http.StatusOK, v)
		}
	}
}

func ListOrders(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, repository.OrdersStruct.Entities)
}

func GetOrderById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.OrdersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			ctx.JSON(http.StatusOK, v)
		}
	}
}

func DeleteOrder(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for i, v := range repository.OrdersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			repository.OrdersStruct.Entities = slices.Delete(repository.OrdersStruct.Entities, i, i+1)
			repository.OrdersStruct.EntitiesLen = len(repository.OrdersStruct.Entities)
			repository.OrdersStruct.SaveToFile("./assets/orders.json")
			ctx.JSON(http.StatusOK, gin.H{"status": "Success", "message": "Объект успешно удален"})
		}
	}
}
