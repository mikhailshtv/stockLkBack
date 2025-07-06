package model

import (
	"time"
)

type OrderStatus int32

const (
	Active OrderStatus = iota + 1
	Executed
	Deleted
)

// сделал конкретные int32 и int64 из-за protobuf и линтера (gosec).
type Order struct {
	ID               int32       `json:"id" bson:"_id,omitempty"`
	Number           int32       `json:"number" bson:"number"`
	TotalCost        int32       `json:"totalCost" bson:"totalCost"`
	CreatedDate      time.Time   `json:"createdDate" bson:"createdDate"`
	LastModifiedDate time.Time   `json:"lastModifiedDate" bson:"lastModifiedDate"`
	Status           OrderStatus `json:"status" bson:"status"`
	Products         []Product   `json:"products" binding:"required" bson:"products"`
}

type OrderRequestBody struct {
	Products []Product `json:"products" binding:"required" bson:"products"`
}
