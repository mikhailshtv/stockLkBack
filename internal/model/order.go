package model

import (
	"time"
)

type OrderStatus int

const (
	Active OrderStatus = iota
	Executed
	Deleted
)

type Order struct {
	Id               uint32
	Number           int
	TotalCost        int
	CreatedDate      time.Time
	lastModifiedDate time.Time
	Status           OrderStatus
}

func (order *Order) LastModifiedDate() time.Time {
	return order.lastModifiedDate
}

func (order *Order) SetLastModifiedDate(lastModifiedDate time.Time) {
	order.lastModifiedDate = lastModifiedDate
}
