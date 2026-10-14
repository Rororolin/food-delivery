package port

import "food-delivery/internal/domain"

// OrderRepository — порт для хранилища заказов
type OrderRepository interface {
	SaveOrder(order *domain.Order) error
	GetOrderByID(orderID string) (*domain.Order, error)
}

// OrderUseCase — порт для сценариев заказа
type OrderUseCase interface {
	PlaceOrder(userID string, items []domain.OrderItem, address, paymentMethod string) (*domain.Order, error)
	TrackOrder(orderID string) (*domain.Order, error)
}
