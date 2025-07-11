package service

import (
	"context"
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
	ctx  context.Context
}

func NewOrdersService(repo repository.Order, ctx context.Context) *OrdersService {
	return &OrdersService{repo: repo, ctx: ctx}
}

func (s *OrdersService) Create(order model.OrderRequestBody, userID int32) (*model.Order, error) {
	createdOrder, err := s.repo.Create(order, userID, s.ctx)
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

func (s *OrdersService) GetAll(userID int32, role model.UserRole) ([]model.Order, error) {
	return s.repo.GetAll(userID, role, s.ctx)
}

func (s *OrdersService) GetByID(id, userID int32, role model.UserRole) (*model.Order, error) {
	return s.repo.GetByID(id, userID, role, s.ctx)
}

func (s *OrdersService) Delete(id int32, userID int32) error {
	delitedOrder, err := s.repo.Delete(id, userID, s.ctx)
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

func (s *OrdersService) Update(id int32, order model.OrderRequestBody, userID int32) (*model.Order, error) {
	updatedOrder, err := s.repo.Update(id, order, userID, s.ctx)
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
