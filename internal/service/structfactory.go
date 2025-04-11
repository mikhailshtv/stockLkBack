package service

import (
	"errors"
	"golang/stockLkBack/internal/model"
	"time"

	"github.com/ddosify/go-faker/faker"
)

func NewEntity() any {
	caseNumber := faker.NewFaker().RandomIntBetween(0, 1)
	switch caseNumber {
	case 0:
		return NewOrder()
	case 1:
		return NewProduct()
	default:
		return errors.New("ошибка при создании сущности")
	}
}

func NewOrder() model.Order {
	order := model.Order{
		Id:          faker.NewFaker().RandomUUID().ID(),
		Number:      faker.NewFaker().RandomIntBetween(1, 999999),
		TotalCost:   faker.NewFaker().RandomIntBetween(1, 999999),
		CreatedDate: time.Now().Add(time.Hour * time.Duration(faker.NewFaker().RandomIntBetween(-8760, -1))),
		Status:      model.OrderStatus(faker.NewFaker().RandomIntBetween(1, 2)),
	}
	order.SetLastModifiedDate(time.Now().Add(time.Hour * time.Duration(faker.NewFaker().RandomIntBetween(-8760, -1))))
	return order
}

func NewProduct() model.Product {
	product := model.Product{
		Id:        faker.NewFaker().RandomUUID().ID(),
		Code:      faker.NewFaker().RandomIntBetween(1, 999999),
		Quantity:  faker.NewFaker().RandomIntBetween(1, 9999),
		Name:      faker.NewFaker().RandomProduct(),
		SalePrice: faker.NewFaker().RandomIntBetween(1, 999999),
	}
	product.SetPurchasePrice(faker.NewFaker().RandomIntBetween(1, 999999))
	return product
}
