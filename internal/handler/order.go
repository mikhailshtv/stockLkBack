package handler

import (
	"golang/stockLkBack/internal/model"
	"log"
	"net/http"
	"strconv"

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
func (h *Handler) CreateOrder(ctx *gin.Context) {
	var order model.OrderRequestBody
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Services.Order.Create(order)
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
func (h *Handler) EditOrder(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	var order model.OrderRequestBody
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderResult, err := h.Services.Order.Update(id, order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orderResult)
}

// OrderList
// @Summary Список заказов
// @Tags Orders
// @Produce		json
// @Success 200 {object} []model.Order
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/orders [get]
// @Security BearerAuth
func (h *Handler) ListOrders(ctx *gin.Context) {
	orders, err := h.Services.Order.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
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
func (h *Handler) GetOrderById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	order, err := h.Services.Order.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
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
func (h *Handler) DeleteOrder(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	err = h.Services.Order.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	success := model.Success{
		Status:  "Success",
		Message: "Объект успешно удален",
	}
	ctx.JSON(http.StatusOK, success)
}
