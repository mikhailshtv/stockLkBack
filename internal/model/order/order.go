package order

import "time"

type Order struct {
	Id               int
	Number           int
	CreatedDate      time.Time
	lastModifiedDate time.Time
	Status           string
	TotalCost        string
}

func (order *Order) GetLastModifiedDate() time.Time {
	return order.lastModifiedDate
}

func (order *Order) SetLastModifiedDate(lastModifiedDate time.Time) {
	order.lastModifiedDate = lastModifiedDate
}
