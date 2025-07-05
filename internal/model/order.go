package model

import (
	"time"
)

type OrderStatus int

const (
	Active OrderStatus = iota + 1
	Executed
	Deleted
)

type Order struct {
	Id               int         `json:"id" bson:"_id,omitempty"`
	Number           int         `json:"number" bson:"number"`
	TotalCost        int         `json:"totalCost" bson:"totalCost"`
	CreatedDate      time.Time   `json:"createdDate" bson:"createdDate"`
	LastModifiedDate time.Time   `json:"lastModifiedDate" bson:"lastModifiedDate"`
	Status           OrderStatus `json:"status" bson:"status"`
	Products         []Product   `json:"products" binding:"required" bson:"products"`
}

type OrderRequestBody struct {
	Products []Product `json:"products" binding:"required" bson:"products"`
}
