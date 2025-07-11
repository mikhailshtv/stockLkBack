package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const NotFoundErrorMessage = "элемент не найден"

//go:generate mockgen -source=repository.go -destination=mocks/repository.go -package=mocks

type Order interface {
	Create(ctx context.Context, order model.OrderRequestBody, userID int32) (*model.Order, error)
	GetAll(ctx context.Context, userID int32, role model.UserRole) ([]model.Order, error)
	GetByID(ctx context.Context, id, userID int32, role model.UserRole) (*model.Order, error)
	Delete(ctx context.Context, id, userID int32) (*model.Order, error)
	Update(ctx context.Context, id int32, orderReq model.OrderRequestBody, userID int32) (*model.Order, error)
	WriteLog(result any, operation, status, tableName string) (int64, error)
}

type Product interface {
	Create(ctx context.Context, product model.Product) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id int32) (*model.Product, error)
	Delete(ctx context.Context, id int32) (*model.Product, error)
	Update(ctx context.Context, id int32, product model.Product) (*model.Product, error)
	WriteLog(result any, operation, status, tableName string) (int64, error)
}

type User interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	Delete(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id int, user model.UserEditBody) (*model.User, error)
	Login(ctx context.Context, user model.LoginRequest) (*model.TokenSuccess, error)
	ChangeUserRole(ctx context.Context, id int, userRoleReq model.UserRoleBody) (*model.User, error)
	ChangePassword(ctx context.Context, id int, changePassworReq model.UserChangePasswordBody) (*model.Success, error)
	WriteLog(result any, operation, status, tableName string) (int64, error)
}

type Repository struct {
	Order
	Product
	User
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		Order:   NewOrdersRepository(db, redis, "ordersCollection"),
		Product: NewProductsRepository(db, redis),
		User:    NewUsersRepository(db, redis),
	}
}

func WriteLog(result any, operation, status, tableName string, redis *redis.Client) (int64, error) {
	incrStr := fmt.Sprintf("%s:id", tableName)
	id, err := redis.Incr(context.TODO(), incrStr).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing ID: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return 0, fmt.Errorf("error marshaling result: %w", err)
	}

	key := fmt.Sprintf("%s:%d", tableName, id)
	_, err = redis.HSet(context.TODO(), key,
		"id", id,
		"operation", operation,
		"status", status,
		"result", resultJSON,
		"date", time.Now().UTC(),
	).Result()
	if err != nil {
		return 0, fmt.Errorf("error saving log: %w", err)
	}
	_, err = redis.Expire(context.TODO(), key, time.Hour*24).Result()
	if err != nil {
		return 0, fmt.Errorf("error setting TTL: %w", err)
	}

	return id, nil
}

func isDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505" // Код ошибки уникальности (unique_violation)
	}
	return false
}
