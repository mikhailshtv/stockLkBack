package service

import (
	"context"
	"log"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
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

func NewOrdersService(ctx context.Context, repo repository.Order) *OrdersService {
	return &OrdersService{repo: repo, ctx: ctx}
}

func (s *OrdersService) Create(order model.OrderRequestBody, userID int) (*model.Order, error) {
	createdOrder, err := s.repo.Create(s.ctx, order, userID)
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

func (s *OrdersService) GetAll(userID int, role model.UserRole) ([]model.Order, error) {
	return s.repo.GetAll(s.ctx, userID, role)
}

func (s *OrdersService) GetByID(id, userID int, role model.UserRole) (*model.Order, error) {
	return s.repo.GetByID(s.ctx, id, userID, role)
}

func (s *OrdersService) Delete(id, userID int) error {
	delitedOrder, err := s.repo.Delete(s.ctx, id, userID)
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

func (s *OrdersService) Update(id int, order model.OrderRequestBody, userID int) (*model.Order, error) {
	updatedOrder, err := s.repo.Update(s.ctx, id, order, userID)
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
