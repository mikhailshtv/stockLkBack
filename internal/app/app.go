package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/config"
	"github.com/mikhailshtv/stockLkBack/internal/handler"
	"github.com/mikhailshtv/stockLkBack/internal/middleware"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg     *config.Config
	ctx     context.Context
	handler *handler.Handler
}

func NewApp(ctx context.Context, cfg *config.Config, handler *handler.Handler) (*App, error) {
	return &App{
		ctx:     ctx,
		cfg:     cfg,
		handler: handler,
	}, nil
}

func (a *App) Start(r *gin.Engine) error {
	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	api := r.Group(a.cfg.Group)
	api.POST("/login", a.handler.Login)
	{
		orders := api.Group("/orders")
		{
			orders.POST("", middleware.TokenAuthMiddleware(), a.handler.CreateOrder)
			orders.PUT("/:id", middleware.TokenAuthMiddleware(), a.handler.EditOrder)
			orders.GET("", middleware.TokenAuthMiddleware(), a.handler.ListOrders)
			orders.GET("/:id", middleware.TokenAuthMiddleware(), a.handler.GetOrderByID)
			orders.DELETE("/:id", middleware.TokenAuthMiddleware(), a.handler.DeleteOrder)
		}
		products := api.Group("/products")
		{
			products.POST("", middleware.TokenAuthMiddleware(), a.handler.CreateProduct)
			products.PUT("/:id", middleware.TokenAuthMiddleware(), a.handler.EditProduct)
			products.GET("", middleware.TokenAuthMiddleware(), a.handler.ListProduct)
			products.GET("/:id", middleware.TokenAuthMiddleware(), a.handler.GetProductByID)
			products.DELETE("/:id", middleware.TokenAuthMiddleware(), a.handler.DeleteProduct)
		}
		users := api.Group("/users")
		{
			users.POST("", a.handler.CreateUser) // фактически регистрация пользователя
			users.PUT("/:id", middleware.TokenAuthMiddleware(), a.handler.EditUser)
			users.GET("", middleware.TokenAuthMiddleware(), a.handler.ListUsers)
			users.GET("/:id", middleware.TokenAuthMiddleware(), a.handler.GetUserByID)
			users.DELETE("/:id", middleware.TokenAuthMiddleware(), a.handler.DeleteUser)
			users.PATCH("/:id/role", middleware.TokenAuthMiddleware(), a.handler.ChangeUserRole)
			users.PATCH("/:id/password", middleware.TokenAuthMiddleware(), a.handler.ChangeUserPassword)
		}
	}

	serverHTTP := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		log.Println("server starting at ", serverHTTP.Addr)
		if err := serverHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := serverHTTP.Shutdown(ctxT); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	log.Println("FINAL server shutdown")
	return nil
}
