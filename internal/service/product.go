package service

import (
	"context"
	"github.com/mikhailshtv/stockLkBack/pkg/errors"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"go.uber.org/zap"
)

const (
	logProductsTableName = "logProduct"
)

type ProductsService struct {
	repo repository.Product
	ctx  context.Context
}

func NewProductsService(ctx context.Context, repo repository.Product) *ProductsService {
	return &ProductsService{repo: repo, ctx: ctx}
}

func (s *ProductsService) Create(product model.Product) (*model.Product, error) {
	createdProduct, err := s.repo.Create(s.ctx, product)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to create product in repository",
			zap.Error(err),
			zap.String("product_name", product.Name),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("product created successfully",
			zap.Int32("product_id", createdProduct.ID),
			zap.String("product_name", createdProduct.Name),
		)
		result = createdProduct
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Create", status, logProductsTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for product creation",
			zap.Error(logErr),
		)
	}
	return createdProduct, err
}

func (s *ProductsService) GetAll(params model.ProductQueryParams) ([]model.Product, error) {
	products, err := s.repo.GetAll(s.ctx, params)
	if err != nil {
		logger.GetLogger().Error("failed to get products from repository",
			zap.Error(err),
		)
		return nil, errors.NewDatabaseError("ошибка получения списка продуктов", err)
	}
	return products, nil
}

func (s *ProductsService) GetTotalCount(params model.ProductQueryParams) (int, error) {
	count, err := s.repo.GetTotalCount(s.ctx, params)
	if err != nil {
		logger.GetLogger().Error("failed to get products count from repository",
			zap.Error(err),
		)
		return 0, errors.NewDatabaseError("ошибка получения количества продуктов", err)
	}
	return count, nil
}

func (s *ProductsService) GetByID(id int) (*model.Product, error) {
	product, err := s.repo.GetByID(s.ctx, id)
	if err != nil {
		logger.GetLogger().Error("failed to get product by ID from repository",
			zap.Error(err),
			zap.Int("product_id", id),
		)
		if err.Error() == "продукт не найден" {
			return nil, errors.NewNotFoundError("продукт", err)
		}
		return nil, errors.NewDatabaseError("ошибка получения продукта", err)
	}
	return product, nil
}

func (s *ProductsService) Delete(id int) error {
	deletedProduct, err := s.repo.Delete(s.ctx, id)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to delete product from repository",
			zap.Error(err),
			zap.Int("product_id", id),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("product deleted successfully",
			zap.Int("product_id", id),
		)
		result = deletedProduct
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Delete", status, logProductsTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for product deletion",
			zap.Error(logErr),
		)
	}
	return err
}

func (s *ProductsService) Update(id int, product model.Product) (*model.Product, error) {
	updatedProduct, err := s.repo.Update(s.ctx, id, product)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to update product in repository",
			zap.Error(err),
			zap.Int("product_id", id),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("product updated successfully",
			zap.Int("product_id", id),
			zap.String("product_name", updatedProduct.Name),
		)
		result = updatedProduct
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Update", status, logProductsTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for product update",
			zap.Error(logErr),
		)
	}
	return updatedProduct, err
}
