package service

import (
	"log"

	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
)

const (
	logErrorStatus     = "Error"
	logSuccessStatus   = "Success"
	logOrdersTableName = "logOrder"
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
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = createdOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Create", status, logOrdersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return createdOrder, err
}

func (s *OrdersService) GetAll() ([]model.Order, error) {
	return s.repo.GetAll()
}

func (s *OrdersService) GetByID(id int32) (*model.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrdersService) Delete(id int32) error {
	delitedOrder, err := s.repo.Delete(id)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = delitedOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Delete", status, logOrdersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return err
}

func (s *OrdersService) Update(id int32, order model.OrderRequestBody) (*model.Order, error) {
	updatedOrder, err := s.repo.Update(id, order)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = updatedOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Update", status, logOrdersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return updatedOrder, err
}
