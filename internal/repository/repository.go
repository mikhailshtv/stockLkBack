package repository

import (
	"golang/stockLkBack/internal/model"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go -package=mocks

type Order interface {
	Create(order model.OrderRequestBody) (*model.Order, error)
	GetAll() ([]model.Order, error)
	GetById(id int) (*model.Order, error)
	Delete(id int) (*model.Order, error)
	Update(id int, order model.OrderRequestBody) (*model.Order, error)
	WriteLog(result any, operation, status string) (int64, error)
}

type Product interface {
	Create(product model.ProductRequestBody) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetById(id int) (*model.Product, error)
	Delete(id int) error
	Update(id int, product model.ProductRequestBody) (*model.Product, error)
	RestoreProductsFromFile(path string)
}

type User interface {
	Create(user model.UserCreateBody) (*model.User, error)
	GetAll() ([]model.User, error)
	GetById(id int) (*model.User, error)
	Delete(id int) error
	Update(id int, user model.UserEditBody) (*model.User, error)
	Login(user model.LoginRequest) (*model.TokenSuccess, error)
	ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error)
	ChangePassword(id int, changePassworReq model.UserChangePasswordBody) (*model.Success, error)
	RestoreUsersFromFile(path string)
}

type Repository struct {
	Order
	Product
	User
}

func NewRepository(db *mongo.Database, redis *redis.Client) *Repository {
	return &Repository{
		Order:   NewOrdersRepository(db, redis, "ordersCollection"),
		Product: NewProductsRepository(db),
		User:    NewUsersRepository(db),
	}
}
