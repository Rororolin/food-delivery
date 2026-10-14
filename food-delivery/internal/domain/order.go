package domain

import "time"

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusPreparing  OrderStatus = "preparing"
	OrderStatusCourier    OrderStatus = "courier_assigned"
	OrderStatusInTransit  OrderStatus = "in_transit"
	OrderStatusDelivered  OrderStatus = "delivered"
)

// OrderItem — позиция в заказе
type OrderItem struct {
	DishID   string `json:"dishId"`
	Quantity int    `json:"quantity"`
}

// Order — заказ
type Order struct {
	ID                string      `json:"orderId"`
	UserID            string      `json:"userId"`
	Items             []OrderItem `json:"items"`
	Address           string      `json:"address"`
	PaymentMethod     string      `json:"paymentMethod"`
	Status            OrderStatus `json:"status"`
	EstimatedDelivery time.Time   `json:"estimatedDelivery"`
	CreatedAt         time.Time   `json:"createdAt"`
}
