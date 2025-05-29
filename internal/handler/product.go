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

func ListProduct(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, repository.ProductsStruct.Entities)
}

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
			ctx.JSON(http.StatusOK, gin.H{"status": "Success", "message": "Объект успешно удален"})
		}
	}
}
