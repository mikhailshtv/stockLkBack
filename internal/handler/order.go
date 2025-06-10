package handler

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"golang/stockLkBack/internal/service"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateOrder
// @Summary Создание заказа
// @Tags Orders
// @Accept			json
// @Produce		json
// @Param order body model.OrderRequestBody true "Объект заказа"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.Error "Invalid request"
// @Router /api/v1/orders [post]
// @Security BearerAuth
func CreateOrder(ctx *gin.Context) {
	var order model.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.SetCommonOrderDataOnCreate(&order)
	repository.CheckAndSaveEntity(order)
	ctx.JSON(http.StatusOK, order)
}

// EditOrder
// @Summary Редактирование заказа
// @Tags Orders
// @Accept			json
// @Produce		json
// @Param order body model.OrderRequestBody true "Объект заказа"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.Error "Invalid request"
// @Param id path string true "id заказа"
// @Router /api/v1/orders{id} [put]
// @Security BearerAuth
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

// OrderList
// @Summary Список заказов
// @Tags Orders
// @Produce		json
// @Success 200 {object} []model.Order
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/orders [get]
// @Security BearerAuth
func ListOrders(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, repository.OrdersStruct.Entities)
}

// GetOrderById
// @Summary Получение заказа по id
// @Tags Orders
// @Produce		json
// @Success 200 {object} model.Order
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id заказа"
// @Router /api/v1/orders/{id} [get]
// @Security BearerAuth
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

// DeleteOrder
// @Summary Удаление заказа
// @Tags Orders
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id заказа"
// @Router /api/v1/orders/{id} [delete]
// @Security BearerAuth
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
			success := model.Success{
				Status:  "Success",
				Message: "Объект успешно удален",
			}
			ctx.JSON(http.StatusOK, success)
		}
	}
}
