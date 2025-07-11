package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/mikhailshtv/stockLkBack/docs"
	"github.com/mikhailshtv/stockLkBack/internal/app"
	"github.com/mikhailshtv/stockLkBack/internal/config"
	"github.com/mikhailshtv/stockLkBack/internal/grpc"
	"github.com/mikhailshtv/stockLkBack/internal/handler"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
	"github.com/mikhailshtv/stockLkBack/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swag init -g cmd/main.go команды для генерации сваггера.
// http://localhost:8080/swagger/index.html посмотреть сваггер.

// @title Сервис управления складом
// @version 1
// @description API для сервиса управления товарами на скаладе продавца

// @host localhost:8080/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// C.hello()
	ctx := context.Background()
	appConfig := config.NewConfig()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке .env файла: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	sqlConfig := repository.SQLConfig{
		Host:           dbHost,
		Port:           dbPort,
		User:           dbUser,
		Password:       dbPass,
		Database:       dbName,
		SSLMode:        "disable",
		MaxConnections: 10,
		Timeout:        5,
	}

	dsn := sqlConfig.CreateDsn()

	// Создание базы данных sql
	db, err := repository.NewSqlxConn(ctx, dsn)
	if err != nil {
		fmt.Println("Ошибка подключения к postgreSQL:", err)
		return
	}

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
	services := service.NewService(ctx, repo)
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
