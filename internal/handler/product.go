//nolint:lll
package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mikhailshtv/stockLkBack/internal/middleware"
	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/pkg/errors"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateProduct
// @Summary Создание продукта
// @Tags Products
// @Accept			json
// @Produce		json
// @Param product body model.ProductRequestBody true "Объект продукта"
// @Success 201 {object} model.Product "Created"
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/products [post]
// @Security BearerAuth.
func (h *Handler) CreateProduct(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	userID := ctx.GetInt(userIDKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	isEmployee := role == model.RoleEmployee
	if !isEmployee {
		middleware.HandleError(ctx, errors.NewForbiddenError("Недостаточно прав для выполнения операции", nil))
		return
	}
	var productReq model.Product
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректное тело запроса", err))
		return
	}

	product, err := h.Services.Product.Create(productReq)
	if err != nil {
		logger.GetLogger().Error("failed to create product",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// EditProduct
// @Summary Редактирование продукта
// @Tags Products
// @Accept			json
// @Produce		json
// @Param product body model.ProductRequestBody true "Объект продукта"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [put]
// @Security BearerAuth.
func (h *Handler) EditProduct(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	userID := ctx.GetInt(userIDKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	isEmployee := role == model.RoleEmployee
	if !isEmployee {
		middleware.HandleError(ctx, errors.NewForbiddenError("Недостаточно прав для выполнения операции", nil))
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID продукта", err))
		return
	}
	var productReq model.Product
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректное тело запроса", err))
		return
	}
	product, err := h.Services.Product.Update(id, productReq)
	if err != nil {
		logger.GetLogger().Error("failed to edit product",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "продукт не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("продукт", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// ListProduct
// @Summary Список продуктов
// @Description Получение списка продуктов с возможностью фильтрации, сортировки и пагинации
// @Tags Products
// @Accept json
// @Produce json
// @Param code query integer false "Фильтр по коду продукта"
// @Param quantity query integer false "Фильтр по количеству"
// @Param name query string false "Фильтр по названию (поиск по подстроке)"
// @Param purchase_price query integer false "Фильтр по закупочной цене"
// @Param sell_price query integer false "Фильтр по цене продажи"
// @Param sort_field query string false "Поле для сортировки" Enums(id, code, quantity, name, purchase_price, sell_price)
// @Param sort_order query string false "Направление сортировки" Enums(ASC, DESC) default(ASC)
// @Param page query integer false "Номер страницы" default(1) minimum(1)
// @Param page_size query integer false "Размер страницы" default(25) minimum(1) maximum(100)
// @Success 200 {object} model.ProductListResponse
// @Failure 400 {object} model.Error "Некорректные параметры запроса"
// @Failure 401 {object} model.Error "Неавторизованный доступ"
// @Failure 500 {object} model.Error "Внутренняя ошибка сервера"
// @Router /api/v1/products [get]
// @Security BearerAuth.
func (h *Handler) ListProduct(ctx *gin.Context) {
	var params model.ProductQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректные параметры запроса", err))
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 25
	}

	products, err := h.Services.Product.GetAll(params)
	if err != nil {
		middleware.HandleError(ctx, err)
		return
	}

	total, err := h.Services.Product.GetTotalCount(params)
	if err != nil {
		middleware.HandleError(ctx, err)
		return
	}

	response := model.ProductListResponse{
		Data:     products,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetProductByID
// @Summary Получение продукта по id
// @Tags Products
// @Produce		json
// @Success 200 {object} model.Product
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [get]
// @Security BearerAuth.
func (h *Handler) GetProductByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	userID := ctx.GetInt(userIDKey)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID продукта", err))
		return
	}
	product, err := h.Services.Product.GetByID(id)
	if err != nil {
		logger.GetLogger().Error("failed to get product",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "продукт не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("продукт", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct
// @Summary Удаление продукта
// @Tags Products
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [delete]
// @Security BearerAuth.
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	userID := ctx.GetInt(userIDKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	isEmployee := role == model.RoleEmployee
	if !isEmployee {
		middleware.HandleError(ctx, errors.NewForbiddenError("Недостаточно прав для выполнения операции", nil))
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID продукта", err))
		return
	}

	if err := h.Services.Product.Delete(id); err != nil {
		logger.GetLogger().Error("failed to delete product",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		if strings.Contains(err.Error(), "продукт не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("продукт", err))
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
