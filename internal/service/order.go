package service

import (
	"context"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
	"github.com/mikhailshtv/stockLkBack/pkg/errors"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"go.uber.org/zap"
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
		logger.GetLogger().Error("failed to create order in repository",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("order created successfully",
			zap.Int("order_id", createdOrder.ID),
			zap.Int("user_id", userID),
		)
		result = createdOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Create", status, logOrdersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for order creation",
			zap.Error(logErr),
		)
	}
	return createdOrder, err
}

func (s *OrdersService) GetAll(userID int, role model.UserRole) ([]model.Order, error) {
	orders, err := s.repo.GetAll(s.ctx, userID, role)
	if err != nil {
		logger.GetLogger().Error("failed to get orders from repository",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return nil, errors.NewDatabaseError("ошибка получения списка заказов", err)
	}
	return orders, nil
}

func (s *OrdersService) GetByID(id, userID int, role model.UserRole) (*model.Order, error) {
	order, err := s.repo.GetByID(s.ctx, id, userID, role)
	if err != nil {
		logger.GetLogger().Error("failed to get order by ID from repository",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		if err.Error() == "заказ не найден" {
			return nil, errors.NewNotFoundError("заказ", err)
		}
		return nil, errors.NewDatabaseError("ошибка получения заказа", err)
	}
	return order, nil
}

func (s *OrdersService) Delete(id, userID int) error {
	deletedOrder, err := s.repo.Delete(s.ctx, id, userID)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to delete order from repository",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("order deleted successfully",
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		result = deletedOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Delete", status, logOrdersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for order deletion",
			zap.Error(logErr),
		)
	}
	return err
}

func (s *OrdersService) Update(id int, order model.OrderRequestBody, userID int) (*model.Order, error) {
	updatedOrder, err := s.repo.Update(s.ctx, id, order, userID)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to update order in repository",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("order updated successfully",
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		result = updatedOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Update", status, logOrdersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for order update",
			zap.Error(logErr),
		)
	}
	return updatedOrder, err
}

func (s *OrdersService) UpdateStatus(
	id int,
	orderStatusRequest model.OrderStatusRequest,
	userID int,
) (*model.Order, error) {
	updatedOrder, err := s.repo.UpdateStatus(s.ctx, id, orderStatusRequest, userID)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to update order status in repository",
			zap.Error(err),
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("order status updated successfully",
			zap.Int("order_id", id),
			zap.Int("user_id", userID),
			zap.String("new_status", orderStatusRequest.Status.Key),
		)
		result = updatedOrder
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "UpdateStatus", status, logOrdersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for order status update",
			zap.Error(logErr),
		)
	}
	return updatedOrder, err
}
