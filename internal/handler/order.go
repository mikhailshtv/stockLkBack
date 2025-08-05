package handler

import (
	"github.com/mikhailshtv/stockLkBack/pkg/errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/mikhailshtv/stockLkBack/internal/middleware"
	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateOrder
// @Summary Создание заказа
// @Tags Orders
// @Accept			json
// @Produce		json
// @Param order body model.OrderRequestBody true "Объект заказа"
// @Success 201 {object} model.Order "Created"
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/orders [post]
// @Security BearerAuth.
func (h *Handler) CreateOrder(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	var orderReq model.OrderRequestBody
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректное тело запроса", err))
		return
	}
	order, err := h.Services.Order.Create(orderReq, userID)
	if err != nil {
		logger.GetLogger().Error("failed to create order",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, order)
}

// EditOrder
// @Summary Редактирование заказа
// @Tags Orders
// @Accept			json
// @Produce		json
// @Param id path string true "id заказа"
// @Param order body model.OrderRequestBody true "Объект заказа"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/orders{id} [put]
// @Security BearerAuth.
func (h *Handler) EditOrder(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("некорректный ID заказа", err))
		return
	}
	var order model.OrderRequestBody
	if err := ctx.ShouldBindJSON(&order); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	orderResult, err := h.Services.Order.Update(id, order, userID)
	if err != nil {
		logger.GetLogger().Error("failed to update order",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "заказ не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("заказ", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orderResult)
}

// ListOrders
// @Summary Список заказов
// @Tags Orders
// @Produce		json
// @Success 200 {object} []model.Order
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 500 {string} string "Internal"
// @Router /api/v1/orders [get]
// @Security BearerAuth.
func (h *Handler) ListOrders(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	orders, err := h.Services.Order.GetAll(userID, role.(model.UserRole))
	if err != nil {
		logger.GetLogger().Error("failed to get orders",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// GetOrderByID
// @Summary Получение заказа по id
// @Tags Orders
// @Produce		json
// @Success 200 {object} model.Order
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id заказа"
// @Router /api/v1/orders/{id} [get]
// @Security BearerAuth.
func (h *Handler) GetOrderByID(ctx *gin.Context) {
	userID := ctx.GetInt(userIDKey)
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("некорректный ID заказа", err))
		return
	}
	order, err := h.Services.Order.GetByID(id, userID, role.(model.UserRole))
	if err != nil {
		logger.GetLogger().Error("failed to get order by ID",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "заказ не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("заказ", err))
			return
		}
		middleware.HandleError(ctx, err)
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
// @Failure 401 {object} model.Error "Unauthorized"
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
		middleware.HandleError(ctx, errors.NewValidationError("некорректный ID заказа", err))
		return
	}
	err = h.Services.Order.Delete(id, userID)
	if err != nil {
		logger.GetLogger().Error("failed to delete order",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "заказ не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("заказ", err))
			return
		}
		middleware.HandleError(ctx, err)
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
// @Param id path string true "id заказа"
// @Param order body model.OrderStatusRequest true "Объект статуса заказа"
// @Success 200 {object} model.Order "Ok"
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {string} string "Internal"
// @Router /api/v1/orders/{id} [patch]
// @Security BearerAuth.
func (h *Handler) ChangeOrderStatus(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	userID := ctx.GetInt(userIDKey)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("некорректный ID заказа", err))
		return
	}
	var order model.OrderStatusRequest
	if err := ctx.ShouldBindJSON(&order); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	orderResult, err := h.Services.Order.UpdateStatus(id, order, userID)
	if err != nil {
		logger.GetLogger().Error("failed to update order status",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "заказ не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("заказ", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orderResult)
}
