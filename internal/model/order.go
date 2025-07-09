package model

import (
	"database/sql/driver"
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
	ID               int32       `json:"id" bson:"_id,omitempty" db:"id"`
	Number           int32       `json:"number" bson:"number" db:"order_number"`
	TotalCost        int32       `json:"totalCost" bson:"totalCost" db:"total_cost"`
	CreatedDate      time.Time   `json:"createdDate" bson:"createdDate" db:"created_date"`
	LastModifiedDate time.Time   `json:"lastModifiedDate" bson:"lastModifiedDate" db:"last_midified_date"`
	Status           OrderStatus `json:"status" bson:"status" db:"status"`
	Products         []Product   `json:"products" binding:"required" bson:"products" db:"-"`
}

type OrderRequestBody struct {
	Products []Product `json:"products" binding:"required" bson:"products"`
}

func (s OrderStatus) Value() (driver.Value, error) {
	return s.Key, nil
}
