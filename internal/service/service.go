package service

import (
	"strconv"

	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Order interface {
	Create(order model.OrderRequestBody) (*model.Order, error)
	GetAll() ([]model.Order, error)
	GetByID(id int32) (*model.Order, error)
	Delete(id int32) error
	Update(id int32, order model.OrderRequestBody) (*model.Order, error)
}

type Product interface {
	Create(product model.ProductRequestBody) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetByID(id int32) (*model.Product, error)
	Delete(id int32) error
	Update(id int32, product model.ProductRequestBody) (*model.Product, error)
}

type User interface {
	Create(user model.UserCreateBody) (*model.User, error)
	GetAll() ([]model.User, error)
	GetByID(id int) (*model.User, error)
	Delete(id int) error
	Update(id int, user model.UserEditBody) (*model.User, error)
	Login(user model.LoginRequest) (*model.TokenSuccess, error)
	ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error)
	ChangePassword(id int, changePassworReq model.UserChangePasswordBody) (*model.Success, error)
}

type Service struct {
	Order
	Product
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Order:   NewOrdersService(repo.Order),
		Product: NewProductsService(repo.Product),
		User:    NewUsersService(repo.User),
	}
}

func ParseInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}
