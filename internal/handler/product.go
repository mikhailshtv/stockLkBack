package handler

import (
	"golang/stockLkBack/internal/model"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProduct
// @Summary Создание продукта
// @Tags Products
// @Accept			json
// @Produce		json
// @Param product body model.ProductRequestBody true "Объект продукта"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.Error "Invalid request"
// @Router /api/v1/products [post]
// @Security BearerAuth
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var productReq model.ProductRequestBody
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.Services.Product.Create(productReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [put]
// @Security BearerAuth
func (h *Handler) EditProduct(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	var productReq model.ProductRequestBody
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.Services.Product.Update(id, productReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// ProductList
// @Summary Список продуктов
// @Tags Products
// @Produce		json
// @Success 200 {object} []model.Product
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/products [get]
// @Security BearerAuth
func (h *Handler) ListProduct(ctx *gin.Context) {
	products, err := h.Services.Product.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// GetProductById
// @Summary Получение продукта по id
// @Tags Products
// @Produce		json
// @Success 200 {object} model.Product
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [get]
// @Security BearerAuth
func (h *Handler) GetProductById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	product, err := h.Services.Product.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct
// @Summary Удаление продукта
// @Tags Products
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id продукта"
// @Router /api/v1/products/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := h.Services.Product.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success := model.Success{
		Status:  "Success",
		Message: "Объект успешно удален",
	}
	ctx.JSON(http.StatusOK, success)
}
