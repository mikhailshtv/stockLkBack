package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang/stockLkBack/internal/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const NotFoundErrorMessage = "элемент не найден"

//go:generate mockgen -source=repository.go -destination=mocks/repository.go -package=mocks

type Order interface {
	Create(order model.OrderRequestBody, userID int32, ctx context.Context) (*model.Order, error)
	GetAll(userID int32, role model.UserRole, ctx context.Context) ([]model.Order, error)
	GetByID(id, userID int32, role model.UserRole, ctx context.Context) (*model.Order, error)
	Delete(id int32, userId int32, ctx context.Context) (*model.Order, error)
	Update(id int32, orderReq model.OrderRequestBody, userID int32, ctx context.Context) (*model.Order, error)
	WriteLog(result any, operation, status, tableName string) (int64, error)
}

type Product interface {
	Create(product model.Product, ctx context.Context) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	GetByID(id int32, ctx context.Context) (*model.Product, error)
	Delete(id int32, ctx context.Context) (*model.Product, error)
	Update(id int32, product model.Product, ctx context.Context) (*model.Product, error)
	WriteLog(result any, operation, status, tableName string) (int64, error)
}

type User interface {
	Create(user model.User, ctx context.Context) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(id int, ctx context.Context) (*model.User, error)
	Delete(id int, ctx context.Context) (*model.User, error)
	Update(id int, user model.UserEditBody, ctx context.Context) (*model.User, error)
	Login(user model.LoginRequest, ctx context.Context) (*model.TokenSuccess, error)
	ChangeUserRole(id int, userRoleReq model.UserRoleBody, ctx context.Context) (*model.User, error)
	ChangePassword(id int, changePassworReq model.UserChangePasswordBody, ctx context.Context) (*model.Success, error)
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
