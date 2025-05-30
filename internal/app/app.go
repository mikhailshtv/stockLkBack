package app

import (
	"context"
	"fmt"
	"golang/stockLkBack/internal/config"
	"golang/stockLkBack/internal/handler"
	"golang/stockLkBack/internal/middleware"
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

func (a *App) Start(r *gin.Engine) error {
	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	api := r.Group(a.cfg.Group)
	api.POST("/login", handler.Login)
	{
		orders := api.Group("/orders")
		{
			orders.POST("", middleware.TokenAuthMiddleware(), handler.CreateOrder)
			orders.PUT("/:id", middleware.TokenAuthMiddleware(), handler.EditOrder)
			orders.GET("", middleware.TokenAuthMiddleware(), handler.ListOrders)
			orders.GET("/:id", middleware.TokenAuthMiddleware(), handler.GetOrderById)
			orders.DELETE("/:id", middleware.TokenAuthMiddleware(), handler.DeleteOrder)
		}
		products := api.Group("/products")
		{
			products.POST("", middleware.TokenAuthMiddleware(), handler.CreateProduct)
			products.PUT("/:id", middleware.TokenAuthMiddleware(), handler.EditProduct)
			products.GET("", middleware.TokenAuthMiddleware(), handler.ListProduct)
			products.GET("/:id", middleware.TokenAuthMiddleware(), handler.GetProductById)
			products.DELETE("/:id", middleware.TokenAuthMiddleware(), handler.DeleteProduct)
		}
		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser) // фактически регистрация пользователя
			users.PUT("/:id", middleware.TokenAuthMiddleware(), handler.EditUser)
			users.GET("", middleware.TokenAuthMiddleware(), handler.ListUsers)
			users.GET("/:id", middleware.TokenAuthMiddleware(), handler.GetUserById)
			users.DELETE("/:id", middleware.TokenAuthMiddleware(), handler.DeleteUser)
			users.PATCH("/:id/role", middleware.TokenAuthMiddleware(), handler.ChangeUserRole)
			users.PATCH("/:id/password", middleware.TokenAuthMiddleware(), handler.ChangeUserPassword)
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
