package service

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
)

type OrdersService struct {
	repo repository.Order
}

func NewOrdersService(repo repository.Order) *OrdersService {
	return &OrdersService{repo: repo}
}

func (s *OrdersService) Create(order model.OrderRequestBody) (*model.Order, error) {
	return s.repo.Create(order)
}

func (s *OrdersService) GetAll() ([]model.Order, error) {
	return s.repo.GetAll()
}

func (s *OrdersService) GetById(id int) (*model.Order, error) {
	return s.repo.GetById(id)
}

func (s *OrdersService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *OrdersService) Update(id int, order model.OrderRequestBody) (*model.Order, error) {
	return s.repo.Update(id, order)
}
