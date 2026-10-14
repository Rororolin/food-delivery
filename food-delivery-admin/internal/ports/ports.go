package ports

import "food-delivery-admin/internal/domain"

// OrderRepository — порт для работы с хранилищем заказов
type OrderRepository interface {
	GetAll(status string) ([]*domain.Order, error)
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, newStatus domain.OrderStatus) (*domain.Order, error)
	AssignCourier(orderID, courierID string) (*domain.Order, error)
	Save(order *domain.Order) error
	MarkRefunded(orderID string) error
}

// UserRepository — порт для работы с хранилищем пользователей
type UserRepository interface {
	GetAll() ([]*domain.User, error)
	GetByID(id string) (*domain.User, error)
	Block(id string) (*domain.User, error)
	Unblock(id string) (*domain.User, error)
	Save(user *domain.User) error
}

// PaymentGateway — порт платёжного шлюза
type PaymentGateway interface {
	ProcessRefund(orderID string, amount float64) (*domain.Refund, error)
	GetRefundStatus(orderID string) (*domain.Refund, error)
}

// NotificationService — порт сервиса уведомлений
type NotificationService interface {
	NotifyCustomer(customerID, message string) error
	NotifyCourier(courierID, message string) error
}
