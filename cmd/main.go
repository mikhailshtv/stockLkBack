package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	_ "golang/stockLkBack/docs"
	"golang/stockLkBack/internal/app"
	"golang/stockLkBack/internal/config"
	"golang/stockLkBack/internal/grpc"
	"golang/stockLkBack/internal/handler"
	"golang/stockLkBack/internal/repository"
	"golang/stockLkBack/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Сервис управления складом
// @version 1
// @description API для сервиса управления товарами на скаладе продавца

// @host localhost:8080/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	ctx := context.Background()
	appConfig := config.NewConfig()
	var builder strings.Builder
	builder.WriteString("mongodb://")
	builder.WriteString(appConfig.DBUsername)
	builder.WriteString(":")
	builder.WriteString(appConfig.DBPassword)
	builder.WriteString("@")
	builder.WriteString("localhost:27017")
	// Подключение к MongoDB.
	clientOptions := options.Client().ApplyURI(builder.String())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println(err.Error())
		}
		fmt.Println("Отключено от MongoDB.")
	}()

	// Пинг сервера для проверки соединения к mongodb
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("Подключено к MongoDB!")

	// Создание или переключение на базу данных mongodb
	db := client.Database(appConfig.DBName)

	// Создание клиента Redis
	clientRedis := redis.NewClient(&redis.Options{
		Addr:     "localhost:8081", // Адрес и порт Redis-сервера
		Password: "",               // Пароль (если есть)
		DB:       0,                // Номер базы данных
	})

	// Проверка соединения к redis
	_, err = clientRedis.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Ошибка подключения к Redis:", err)
		return
	}

	repo := repository.NewRepository(db, clientRedis)
	repo.Product.RestoreProductsFromFile("./assets/products.json")
	repo.User.RestoreUsersFromFile("./assets/users.json")
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	go grpc.StartServer(handlers)

	newApp, err := app.NewApp(ctx, appConfig, handlers)
	if err != nil {
		log.Println(err.Error())
		return
	}
	r := gin.Default()
	url := ginSwagger.URL("http://localhost:8080/api/v1/swagger/doc.json")
	r.GET("api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err = newApp.Start(r)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
