package service

import (
	"context"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"log"
)

const (
	logProductsTableName = "logProduct"
)

type ProductsService struct {
	repo repository.Product
	ctx  context.Context
}

func NewProductsService(repo repository.Product, ctx context.Context) *ProductsService {
	return &ProductsService{repo: repo, ctx: ctx}
}

func (s *ProductsService) Create(product model.Product) (*model.Product, error) {
	createdProduct, err := s.repo.Create(product, s.ctx)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = createdProduct
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Create", status, logProductsTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return createdProduct, err
}

func (s *ProductsService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll(s.ctx)
}

func (s *ProductsService) GetByID(id int32) (*model.Product, error) {
	return s.repo.GetByID(id, s.ctx)
}

func (s *ProductsService) Delete(id int32) error {
	delitedProduct, err := s.repo.Delete(id, s.ctx)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = delitedProduct
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Delete", status, logOrdersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return err
}

func (s *ProductsService) Update(id int32, product model.Product) (*model.Product, error) {
	updatedProduct, err := s.repo.Update(id, product, s.ctx)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = updatedProduct
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Update", status, logOrdersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return updatedProduct, err
}
