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
	Id               int         `json:"id"`
	Number           int         `json:"number"`
	TotalCost        int         `json:"totalCost"`
	CreatedDate      time.Time   `json:"createdDate"`
	LastModifiedDate time.Time   `json:"lastModifiedDate"`
	Status           OrderStatus `json:"status"`
	Products         []Product   `json:"products" binding:"required"`
}

// func (order *Order) LastModifiedDate() time.Time {
// 	return order.lastModifiedDate
// }

// func (order *Order) SetLastModifiedDate(lastModifiedDate time.Time) {
// 	order.lastModifiedDate = lastModifiedDate
// }
