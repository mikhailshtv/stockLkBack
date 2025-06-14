package repository

import (
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/model"
	"io"
	"log"
	"os"
	"slices"
)

type ProductsRepository struct {
	Products   []model.Product
	ProductLen int
}

func NewProductsRepository() *ProductsRepository {
	return &ProductsRepository{Products: []model.Product{}}
}

func (pr *ProductsRepository) Create(productRequest model.ProductRequestBody) (*model.Product, error) {
	var product model.Product
	if pr.ProductLen > 0 {
		lastProduct := pr.Products[pr.ProductLen-1]
		product.Id = lastProduct.Id + 1
	} else {
		product.Id = 1
	}
	product.Code = productRequest.Code
	product.Name = productRequest.Name
	product.Quantity = productRequest.Quantity
	product.PurchasePrice = productRequest.PurchasePrice
	product.SalePrice = productRequest.SalePrice

	pr.Products = append(pr.Products, product)
	pr.ProductLen = len(pr.Products)

	if err := saveProductsToFile(pr.Products); err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *ProductsRepository) GetAll() ([]model.Product, error) {
	return pr.Products, nil
}

func (pr *ProductsRepository) GetById(id int) (*model.Product, error) {
	idx := slices.IndexFunc(pr.Products, func(product model.Product) bool { return product.Id == id })
	if idx == -1 {
		return nil, errors.New("элемент не найден")
	}
	return &pr.Products[idx], nil
}

func (pr *ProductsRepository) Delete(id int) error {
	pr.Products = slices.DeleteFunc(pr.Products, func(product model.Product) bool { return product.Id == id })
	pr.ProductLen = len(pr.Products)
	if err := saveProductsToFile(pr.Products); err != nil {
		return err
	}
	return nil
}

func (pr *ProductsRepository) Update(id int, product model.ProductRequestBody) (*model.Product, error) {
	idx := slices.IndexFunc(pr.Products, func(product model.Product) bool { return product.Id == id })
	if idx == -1 {
		return nil, errors.New("элемент не найден")
	}
	pr.Products[idx].Code = product.Code
	pr.Products[idx].Name = product.Name
	pr.Products[idx].Quantity = product.Quantity
	pr.Products[idx].PurchasePrice = product.PurchasePrice
	pr.Products[idx].SalePrice = product.SalePrice

	if err := saveProductsToFile(pr.Products); err != nil {
		return nil, err
	}
	return &pr.Products[idx], nil
}

func (pr *ProductsRepository) RestoreProductsFromFile(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v\n", err.Error())
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Ошибка чтения из файла: %v\n", err.Error())
	}
	if len(data) == 0 {
		return
	}

	jsonError := json.Unmarshal(data, &pr.Products)
	pr.ProductLen = len(pr.Products)
	if jsonError != nil {
		log.Fatalf("Ошибка десериализации: %v\n", jsonError.Error())
	}
}

func saveProductsToFile(products []model.Product) error {
	outputPath := "./assets"
	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Ошибка создания каталога: %v\n", err.Error())
			return err
		}
	}
	path := "./assets/products.json"
	json, err := json.Marshal(products)
	if err != nil {
		log.Fatalf("Ошибка конвертирования в json: %v\n", err.Error())
		return err
	}
	if err := os.WriteFile(path, json, os.ModePerm); err != nil {
		log.Fatalf("Ошибка записи в файл: %v\n", err.Error())
		return err
	}
	return nil
}
