package handler

import (
	"net/http"
	"strings"

	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/service"

	"github.com/gin-gonic/gin"
)

// CreateProduct
// @Summary Создание продукта
// @Tags Products
// @Accept			json
// @Produce		json
// @Param product body model.ProductRequestBody true "Объект продукта"
// @Success 201 {object} model.Product "Created"
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/products [post]
// @Security BearerAuth.
func (h *Handler) CreateProduct(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная роль пользователя"})
		return
	}
	isEmployee := role == "employee"
	if !isEmployee {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно прав для выполнения операции"})
		return
	}
	var productReq model.Product
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.Services.Product.Create(productReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [put]
// @Security BearerAuth.
func (h *Handler) EditProduct(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректная роль пользователя"})
		return
	}
	isEmployee := role == "employee"
	if !isEmployee {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно прав для выполнения операции"})
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := service.ParseInt32(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID продукта"})
		return
	}
	var productReq model.Product
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.Services.Product.Update(id, productReq)
	if err != nil {
		if strings.Contains(err.Error(), "продукт не найден") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// ProductList
// @Summary Список продуктов
// @Tags Products
// @Produce		json
// @Success 200 {object} []model.Product
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/products [get]
// @Security BearerAuth.
func (h *Handler) ListProduct(ctx *gin.Context) {
	products, err := h.Services.Product.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// GetProductById
// @Summary Получение продукта по id
// @Tags Products
// @Produce		json
// @Success 200 {object} model.Product
// @Failure 400 {string} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [get]
// @Security BearerAuth.
func (h *Handler) GetProductByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := service.ParseInt32(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID продукта"})
		return
	}
	product, err := h.Services.Product.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "продукт не найден") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [delete]
// @Security BearerAuth.
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	isEmployee := ctx.GetString(userRoleKey) == "employee"
	if !isEmployee {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно прав для выполнения операции"})
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := service.ParseInt32(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID продукта"})
		return
	}

	if err := h.Services.Product.Delete(id); err != nil {
		if strings.Contains(err.Error(), "продукт не найден") {
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
