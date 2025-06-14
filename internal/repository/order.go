package repository

import (
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/model"
	"io"
	"log"
	"os"
	"slices"
	"time"
)

type OrdersRepository struct {
	Orders    []model.Order
	OrdersLen int
}

func NewOrdersRepository() *OrdersRepository {
	return &OrdersRepository{Orders: []model.Order{}}
}

func (or *OrdersRepository) Create(orderRequest model.OrderRequestBody) (*model.Order, error) {
	var order model.Order
	if or.OrdersLen > 0 {
		lastOrder := or.Orders[or.OrdersLen-1]
		order.Id = lastOrder.Id + 1
		order.Number = lastOrder.Number + 1
	} else {
		order.Id = 1
		order.Number = 1
	}
	order.CreatedDate = time.Now().UTC()
	order.LastModifiedDate = time.Now().UTC()
	order.Status = model.Active
	totalCost := 0
	for _, product := range orderRequest.Products {
		totalCost += product.SalePrice
	}
	order.TotalCost = totalCost
	order.Products = orderRequest.Products
	or.Orders = append(or.Orders, order)
	or.OrdersLen = len(or.Orders)

	if err := saveOrdersToFile(or.Orders); err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *OrdersRepository) GetAll() ([]model.Order, error) {
	return or.Orders, nil
}

func (or *OrdersRepository) GetById(id int) (*model.Order, error) {
	idx := slices.IndexFunc(or.Orders, func(order model.Order) bool { return order.Id == id })
	if idx == -1 {
		return nil, errors.New("элемент не найден")
	}
	return &or.Orders[idx], nil
}

func (or *OrdersRepository) Delete(id int) error {
	or.Orders = slices.DeleteFunc(or.Orders, func(order model.Order) bool { return order.Id == id })
	or.OrdersLen = len(or.Orders)
	if err := saveOrdersToFile(or.Orders); err != nil {
		return err
	}
	return nil
}

func (or *OrdersRepository) Update(id int, order model.OrderRequestBody) (*model.Order, error) {
	idx := slices.IndexFunc(or.Orders, func(order model.Order) bool { return order.Id == id })
	if idx == -1 {
		return nil, errors.New("элемент не найден")
	}
	or.Orders[idx].Products = order.Products
	or.Orders[idx].LastModifiedDate = time.Now().UTC()
	or.Orders[idx].TotalCost = 0
	for _, product := range order.Products {
		or.Orders[idx].TotalCost += product.SalePrice
	}
	if err := saveOrdersToFile(or.Orders); err != nil {
		return nil, err
	}
	return &or.Orders[idx], nil
}

func (or *OrdersRepository) RestoreOrdersFromFile(path string) {
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

	jsonError := json.Unmarshal(data, &or.Orders)
	or.OrdersLen = len(or.Orders)
	if jsonError != nil {
		log.Fatalf("Ошибка десериализации: %v\n", jsonError.Error())
	}
}

func saveOrdersToFile(orders []model.Order) error {
	outputPath := "./assets"
	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Ошибка создания каталога: %v\n", err.Error())
			return err
		}
	}
	path := "./assets/orders.json"
	json, err := json.Marshal(orders)
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
