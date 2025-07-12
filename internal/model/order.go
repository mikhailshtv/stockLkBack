package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type OrderStatus struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
}

var (
	StatusActive = OrderStatus{
		Key:         "active",
		DisplayName: "Активный",
	}

	StatusExecuted = OrderStatus{
		Key:         "executed",
		DisplayName: "Выполнен",
	}

	StatusDeleted = OrderStatus{
		Key:         "deleted",
		DisplayName: "Удален",
	}
)

// сделал конкретные int32 из-за protobuf и линтера (gosec).

type Order struct {
	ID               int         `json:"id" bson:"_id,omitempty" db:"id"`
	Number           int         `json:"number" bson:"number" db:"order_number"`
	TotalCost        int         `json:"totalCost" bson:"totalCost" db:"total_cost"`
	CreatedDate      time.Time   `json:"createdDate" bson:"createdDate" db:"created_date"`
	LastModifiedDate time.Time   `json:"lastModifiedDate" bson:"lastModifiedDate" db:"last_modified_date"`
	Status           OrderStatus `json:"status" bson:"status" db:"status"`
	Products         []Product   `json:"products" binding:"required" bson:"products" db:"-"`
	UserID           int         `json:"userId" db:"user_id"`
}

type OrderRequestBody struct {
	Products []OrderProduct `json:"products" binding:"required" bson:"products"`
}

type OrderProduct struct {
	ProductID int `json:"productId" bindings:"required" db:"product_id"`
	Quantity  int `json:"quantity" bindings:"required" db:"quantity"`    // Количество покупаемых товаров
	SellPrice int `json:"sellPrice" bindings:"required" db:"sell_price"` // Цена товара на момент создания заказа
}

func (os *OrderStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("неверный тип для OrderStatus")
	}

	switch str {
	case "active":
		*os = StatusActive
	case "executed":
		*os = StatusExecuted
	case "deleted":
		*os = StatusDeleted
	default:
		return fmt.Errorf("неизвестный статус: %s", str)
	}

	return nil
}

func (os OrderStatus) Value() (driver.Value, error) {
	return os.Key, nil
}
