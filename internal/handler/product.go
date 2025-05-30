package handler

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"log"
	"net/http"
	"slices"
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
func CreateProduct(ctx *gin.Context) {
	var product model.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Id = repository.ProductsStruct.EntitiesLen + 1
	repository.CheckAndSaveEntity(product)
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
func EditProduct(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.ProductsStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			var product model.Product
			if err := ctx.ShouldBindJSON(&product); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			v.Code = product.Code
			v.Name = product.Name
			v.PurchasePrice = product.PurchasePrice
			v.SalePrice = product.SalePrice
			v.Quantity = product.Quantity
			repository.ProductsStruct.SaveToFile("./assets/products.json")
			ctx.JSON(http.StatusOK, v)
		}
	}
}

// ProductList
// @Summary Список продуктов
// @Tags Products
// @Produce		json
// @Success 200 {object} []model.Product
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/products [get]
// @Security BearerAuth
func ListProduct(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, repository.ProductsStruct.Entities)
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
func GetProductById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.ProductsStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			ctx.JSON(http.StatusOK, v)
		}
	}
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
func DeleteProduct(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for i, v := range repository.ProductsStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			repository.ProductsStruct.Entities = slices.Delete(repository.ProductsStruct.Entities, i, i+1)
			repository.ProductsStruct.EntitiesLen = len(repository.ProductsStruct.Entities)
			repository.ProductsStruct.SaveToFile("./assets/products.json")
			success := model.Success{
				Status:  "Success",
				Message: "Объект успешно удален",
			}
			ctx.JSON(http.StatusOK, success)
		}
	}
}
