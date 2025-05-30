package main

import (
	"context"
	"golang/stockLkBack/internal/app"
	"golang/stockLkBack/internal/config"
	"golang/stockLkBack/internal/service"
	"log"

	_ "golang/stockLkBack/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Сервис управления складом
// @version 1
// @description API для сервиса управления товарами на скаладе продавца

// @host localhost:8080/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	service.RestoreData()
	r := gin.Default()

	newApp, err := app.NewApp(context.Background(), config.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	url := ginSwagger.URL("http://localhost:8080/api/v1/swagger/doc.json")
	r.GET("api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err = newApp.Start(r)
	if err != nil {
		log.Fatal(err)
	}

}
