package service

import (
	"context"
	"strconv"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Order interface {
	Create(order model.OrderRequestBody, userID int) (*model.Order, error)
	GetAll(userID int, role model.UserRole) ([]model.Order, error)
	GetByID(id, userID int, role model.UserRole) (*model.Order, error)
	Delete(id, userID int) error
	Update(id int, order model.OrderRequestBody, userID int) (*model.Order, error)
}

type Product interface {
	Create(product model.Product) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetByID(id int) (*model.Product, error)
	Delete(id int) error
	Update(id int, product model.Product) (*model.Product, error)
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

func NewService(ctx context.Context, repo *repository.Repository) *Service {
	return &Service{
		Order:   NewOrdersService(ctx, repo.Order),
		Product: NewProductsService(ctx, repo.Product),
		User:    NewUsersService(ctx, repo.User),
	}
}

func ParseInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}
