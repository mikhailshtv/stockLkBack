package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mikhailshtv/stockLkBack/internal/model"

	"github.com/gin-gonic/gin"
)

// CreateOrder
// @Summary Создание заказа
// @Tags Orders
// @Accept			json
// @Produce		json
// @Param order body model.OrderRequestBody true "Объект заказа"
// @Success 201 {object} model.Order "Created"
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/orders [post]
// @Security BearerAuth.
func (h *Handler) CreateOrder(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	var orderReq model.OrderRequestBody
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное тело запроса"})
		return
	}
	order, err := h.Services.Order.Create(orderReq, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id заказа"
// @Router /api/v1/orders{id} [put]
// @Security BearerAuth.
func (h *Handler) EditOrder(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID заказа"})
		return
	}
	var order model.OrderRequestBody
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderResult, err := h.Services.Order.Update(id, order, userID)
	if err != nil {
		if strings.Contains(err.Error(), "заказ не найден") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orderResult)
}

// OrderList
// @Summary Список заказов
// @Tags Orders
// @Produce		json
// @Success 200 {object} []model.Order
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 500 {string} string "Internal"
// @Router /api/v1/orders [get]
// @Security BearerAuth.
func (h *Handler) ListOrders(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная роль пользователя"})
		return
	}
	orders, err := h.Services.Order.GetAll(userID, role.(model.UserRole))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// GetOrderById
// @Summary Получение заказа по id
// @Tags Orders
// @Produce		json
// @Success 200 {object} model.Order
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id заказа"
// @Router /api/v1/orders/{id} [get]
// @Security BearerAuth.
func (h *Handler) GetOrderByID(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная роль пользователя"})
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID заказа"})
		return
	}
	order, err := h.Services.Order.GetByID(id, userID, role.(model.UserRole))
	if err != nil {
		if strings.Contains(err.Error(), "заказ не найден") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

// DeleteOrder
// @Summary Удаление заказа
// @Tags Orders
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {string} string "Internal"
// @Param id path string true "id заказа"
// @Router /api/v1/orders/{id} [delete]
// @Security BearerAuth.
func (h *Handler) DeleteOrder(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	userID := ctx.GetInt(userIDKey)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID заказа"})
		return
	}
	err = h.Services.Order.Delete(id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "заказ не найден") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success := model.Success{
		Status:  "Success",
		Message: "Объект успешно удален",
	}
	ctx.JSON(http.StatusOK, success)
}

// ChangeOrderStatus
// @Summary Изменение статуса заказа
// @Tags Orders
// @Produce		json
// @Success 200 {object} model.Order "Ok"
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {string} string "Internal"
// @Param id path string true "id заказа"
// @Router /api/v1/orders/{id} [delete]
// @Security BearerAuth.
func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	userID := ctx.GetInt(userIDKey)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID заказа"})
		return
	}
	var order model.OrderStatusRequest
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderResult, err := h.Services.Order.UpdateStatus(id, order, userID)
	if err != nil {
		if strings.Contains(err.Error(), "заказ не найден") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orderResult)
}
