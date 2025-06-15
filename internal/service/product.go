package service

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
)

type ProductsService struct {
	repo repository.Product
}

func NewProductsService(repo repository.Product) *ProductsService {
	return &ProductsService{repo: repo}
}

func (s *ProductsService) Create(product model.ProductRequestBody) (*model.Product, error) {
	return s.repo.Create(product)
}

func (s *ProductsService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductsService) GetById(id int) (*model.Product, error) {
	return s.repo.GetById(id)
}

func (s *ProductsService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductsService) Update(id int, product model.ProductRequestBody) (*model.Product, error) {
	return s.repo.Update(id, product)
}
