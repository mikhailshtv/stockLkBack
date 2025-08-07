package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/mikhailshtv/stockLkBack/config"
	_ "github.com/mikhailshtv/stockLkBack/docs"
	"github.com/mikhailshtv/stockLkBack/internal/app"
	"github.com/mikhailshtv/stockLkBack/internal/grpc"
	"github.com/mikhailshtv/stockLkBack/internal/handler"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
	"github.com/mikhailshtv/stockLkBack/internal/service"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// swag init -g cmd/main.go команды для генерации сваггера.
// http://localhost:8080/swagger/index.html посмотреть сваггер.

// @title Сервис управления складом
// @version 1
// @description API для сервиса управления товарами на складе продавца

// @host localhost:8080/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	ctx := context.Background()
	var configPath string

	flag.StringVar(&configPath, "c", "config/config.yaml", "path to config")
	flag.Parse()

	cfg, err := config.Read(configPath)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	if err := logger.InitLogger(cfg.Logging.Level); err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}
	defer logger.Sync()

	logger.GetLogger().Info("application starting",
		zap.String("config_path", configPath),
		zap.String("log_level", cfg.Logging.Level),
	)

	sqlConfig := repository.SQLConfig{
		Host:           cfg.DB.Host,
		Port:           cfg.DB.Port,
		User:           os.Getenv(cfg.DB.UserEnvKey),
		Password:       os.Getenv(cfg.DB.PassEnvKey),
		Database:       cfg.DB.DBName,
		SSLMode:        "disable",
		MaxConnections: cfg.DB.MaxConnections,
		Timeout:        cfg.DB.Timeout,
	}

	dsn := sqlConfig.CreateDsn()

	db, err := repository.NewSqlxConn(ctx, dsn)
	if err != nil {
		logger.GetLogger().Error("failed to connect to database", zap.Error(err))
		return
	}

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: "",
		DB:       cfg.Redis.DB,
	})

	_, err = clientRedis.Ping(ctx).Result()
	if err != nil {
		logger.GetLogger().Error("failed to connect to redis", zap.Error(err))
		return
	}

	repo := repository.NewRepository(db, clientRedis)
	services := service.NewService(ctx, repo)
	handlers := handler.NewHandler(services)

	go grpc.StartServer(handlers)

	newApp, err := app.NewApp(ctx, cfg, handlers)
	if err != nil {
		log.Println(err.Error())
		return
	}
	r := gin.Default()
	url := ginSwagger.URL("/api/v1/swagger/doc.json")
	r.GET("api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	err = newApp.Start(r)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
