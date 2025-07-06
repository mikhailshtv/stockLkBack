package repository

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"slices"

	"golang/stockLkBack/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProductsRepository struct {
	Products   []model.Product
	ProductLen int
	db         *mongo.Database
}

func NewProductsRepository(db *mongo.Database) *ProductsRepository {
	return &ProductsRepository{db: db}
}

func (pr *ProductsRepository) Create(productRequest model.ProductRequestBody) (*model.Product, error) {
	var product model.Product
	if pr.ProductLen > 0 {
		lastProduct := pr.Products[pr.ProductLen-1]
		product.ID = lastProduct.ID + 1
	} else {
		product.ID = 1
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

func (pr *ProductsRepository) GetByID(id int32) (*model.Product, error) {
	idx := slices.IndexFunc(pr.Products, func(product model.Product) bool { return product.ID == id })
	if idx == -1 {
		return nil, errors.New(NotFoundErrorMessage)
	}
	return &pr.Products[idx], nil
}

func (pr *ProductsRepository) Delete(id int32) error {
	pr.Products = slices.DeleteFunc(pr.Products, func(product model.Product) bool { return product.ID == id })
	pr.ProductLen = len(pr.Products)
	if err := saveProductsToFile(pr.Products); err != nil {
		return err
	}
	return nil
}

func (pr *ProductsRepository) Update(id int32, product model.ProductRequestBody) (*model.Product, error) {
	idx := slices.IndexFunc(pr.Products, func(product model.Product) bool { return product.ID == id })
	if idx == -1 {
		return nil, errors.New(NotFoundErrorMessage)
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
		log.Printf("Ошибка открытия файла: %v\n", err.Error())
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка чтения из файла: %v\n", err.Error())
		return
	}
	if len(data) == 0 {
		return
	}

	jsonError := json.Unmarshal(data, &pr.Products)
	pr.ProductLen = len(pr.Products)
	if jsonError != nil {
		log.Printf("Ошибка десериализации: %v\n", jsonError.Error())
		return
	}
}

func saveProductsToFile(products []model.Product) error {
	outputPath := "./assets"
	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputPath, os.ModePerm)
		if err != nil {
			log.Printf("Ошибка создания каталога: %v\n", err.Error())
			return err
		}
	}
	path := "./assets/products.json"
	json, err := json.Marshal(products)
	if err != nil {
		log.Printf("Ошибка конвертирования в json: %v\n", err.Error())
		return err
	}
	if err := os.WriteFile(path, json, 0o600); err != nil {
		log.Printf("Ошибка записи в файл: %v\n", err.Error())
		return err
	}
	return nil
}
