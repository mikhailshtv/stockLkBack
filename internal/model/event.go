package model

import "time"

type OrderEventType string

const (
	OrderEventCreated       OrderEventType = "ORDER_CREATED"
	OrderEventUpdated       OrderEventType = "ORDER_UPDATED"
	OrderEventStatusChanged OrderEventType = "ORDER_STATUS_CHANGED"
)

type OrderEvent struct {
	EventType   OrderEventType `json:"eventType"`
	OrderID     int            `json:"orderId"`
	OrderNumber int            `json:"orderNumber"`
	TotalCost   int            `json:"totalCost"`
	Status      string         `json:"status"`
	UserID      int            `json:"userId"`
	UserEmail   string         `json:"userEmail"`
	UserName    string         `json:"userName"`
	OccurredAt  time.Time      `json:"occurredAt"`
}
