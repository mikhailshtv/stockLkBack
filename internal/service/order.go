package service

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"log"
)

type OrdersService struct {
	repo repository.Order
}

func NewOrdersService(repo repository.Order) *OrdersService {
	return &OrdersService{repo: repo}
}

func (s *OrdersService) Create(order model.OrderRequestBody) (*model.Order, error) {
	createdOrder, err := s.repo.Create(order)
	var result any
	var status string
	if err != nil {
		result = err
		status = "Error"
	} else {
		result = createdOrder
		status = "Success"
	}

	_, logErr := s.repo.WriteLog(result, "Create", status)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return createdOrder, err
}

func (s *OrdersService) GetAll() ([]model.Order, error) {
	return s.repo.GetAll()
}

func (s *OrdersService) GetById(id int) (*model.Order, error) {
	return s.repo.GetById(id)
}

func (s *OrdersService) Delete(id int) error {
	delitedOrder, err := s.repo.Delete(id)
	var result any
	var status string
	if err != nil {
		result = err
		status = "Error"
	}

	result = delitedOrder
	status = "Success"

	_, logErr := s.repo.WriteLog(result, "Delete", status)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return err
}

func (s *OrdersService) Update(id int, order model.OrderRequestBody) (*model.Order, error) {
	updatedOrder, err := s.repo.Update(id, order)
	var result any
	var status string
	if err != nil {
		result = err
		status = "Error"
	}
	result = updatedOrder
	status = "Success"

	_, logErr := s.repo.WriteLog(result, "Update", status)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return updatedOrder, err
}
