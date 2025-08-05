package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikhailshtv/stockLkBack/config"
	"github.com/mikhailshtv/stockLkBack/internal/handler"
	"github.com/mikhailshtv/stockLkBack/internal/middleware"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"go.uber.org/zap"
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

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())

	api := r.Group(a.cfg.HTTP.BasePath)
	api.POST("/login", a.handler.Login)
	{
		orders := api.Group("/orders")
		{
			orders.POST("", middleware.TokenAuthMiddleware(), a.handler.CreateOrder)
			orders.PUT("/:id", middleware.TokenAuthMiddleware(), a.handler.EditOrder)
			orders.GET("", middleware.TokenAuthMiddleware(), a.handler.ListOrders)
			orders.GET("/:id", middleware.TokenAuthMiddleware(), a.handler.GetOrderByID)
			orders.DELETE("/:id", middleware.TokenAuthMiddleware(), a.handler.DeleteOrder)
			orders.PATCH("/:id", middleware.TokenAuthMiddleware(), a.handler.ChangeOrderStatus)
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
		Addr:              a.cfg.HTTP.Address,
		Handler:           useCors(r, a.cfg.HTTP.AllowedOrigins),
		ReadHeaderTimeout: a.cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       a.cfg.HTTP.ReadTimeout,
		WriteTimeout:      a.cfg.HTTP.WriteTimeout,
		IdleTimeout:       a.cfg.HTTP.IdleTimeout,
	}

	go func() {
		logger.GetLogger().Info("server starting",
			zap.String("address", serverHTTP.Addr),
		)
		if err := serverHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Fatal("server failed to start",
				zap.Error(err),
			)
		}
	}()

	<-ctx.Done()
	logger.GetLogger().Info("got interruption signal")
	ctxT, cancel := context.WithTimeout(ctx, a.cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := serverHTTP.Shutdown(ctxT); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}
	logger.GetLogger().Info("server shutdown completed")
	return nil
}

func useCors(h http.Handler, allowedOrigins []string) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return corsMiddleware.Handler(h)
}
