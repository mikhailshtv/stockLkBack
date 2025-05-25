package handler

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"net/http"

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
