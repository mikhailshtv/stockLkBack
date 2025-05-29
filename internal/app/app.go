package app

import (
	"context"
	"fmt"
	"golang/stockLkBack/internal/config"
	"golang/stockLkBack/internal/handler"
	"golang/stockLkBack/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg *config.Config
	ctx context.Context
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	return &App{
		ctx: ctx,
		cfg: cfg,
	}, nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group(a.cfg.Group)
	{
		orders := api.Group("/orders")
		{
			orders.POST("", handler.CreateOrder)
			orders.PUT("/:id", handler.EditOrder)
			orders.GET("", handler.ListOrders)
			orders.GET("/:id", handler.GetOrderById)
			orders.DELETE("/:id", handler.DeleteOrder)
		}
		products := api.Group("/products")
		{
			products.POST("", handler.CreateProduct)
			products.PUT("/:id", handler.EditProduct)
			products.GET("", handler.ListProduct)
			products.GET("/:id", handler.GetProductById)
			products.DELETE("/:id", handler.DeleteProduct)
		}
		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser)
			users.PUT("/:id", handler.EditUser)
			users.GET("", handler.ListUsers)
			users.GET("/:id", handler.GetUserById)
			users.DELETE("/:id", handler.DeleteUser)
		}
	}

	serverHTTP := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port),
		Handler: r,
	}

	go func() {
		log.Println("server starting at ", serverHTTP.Addr)
		if err := serverHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	service.LogAddedEntities(ctx)

	<-ctx.Done()
	log.Println("got interruption signal")
	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := serverHTTP.Shutdown(ctxT); err != nil {
		return fmt.Errorf("shutdown server: %s", err)
	}
	log.Println("FINAL server shutdown")
	return nil
}
