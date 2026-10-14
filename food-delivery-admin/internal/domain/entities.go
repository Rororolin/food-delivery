package domain

import "time"

// OrderStatus — возможные статусы заказа
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAssigned   OrderStatus = "assigned"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// ValidTransitions — допустимые переходы статусов
var ValidTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:    {OrderStatusAssigned, OrderStatusDelivering, OrderStatusDelivered, OrderStatusCancelled},
	OrderStatusAssigned:   {OrderStatusDelivering, OrderStatusDelivered, OrderStatusCancelled},
	OrderStatusDelivering: {OrderStatusDelivered, OrderStatusCancelled},
	OrderStatusDelivered:  {},
	OrderStatusCancelled:  {},
}

// Order — доменная сущность «Заказ»
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	CourierID   string      `json:"courier_id,omitempty"`
	Status      OrderStatus `json:"status"`
	TotalAmount float64     `json:"totalAmount"`
	IsPaid      bool        `json:"is_paid"`
	IsRefunded  bool        `json:"is_refunded"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// CanTransitionTo — проверяет допустимость перехода статуса
func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	allowed, ok := ValidTransitions[o.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

// UserRole — роль пользователя в системе
type UserRole string

const (
	RoleCustomer  UserRole = "customer"
	RoleCourier   UserRole = "courier"
	RoleAdmin     UserRole = "admin"
)

// User — доменная сущность «Пользователь»
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"createdAt"`
}

// RefundStatus — статус возврата
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusFailed    RefundStatus = "failed"
)

// Refund — доменная сущность «Возврат средств»
type Refund struct {
	OrderID   string       `json:"order_id"`
	Amount    float64      `json:"amount"`
	Status    RefundStatus `json:"refund_status"`
	CreatedAt time.Time    `json:"createdAt"`
}
