package repository

import (
	"errors"
	"food-delivery/internal/domain"
	"sync"
)

// InMemoryOrderRepository — хранилище заказов в памяти
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

// SaveOrder сохраняет заказ
func (r *InMemoryOrderRepository) SaveOrder(order *domain.Order) error {
	if order == nil {
		return errors.New("order must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

// GetOrderByID возвращает заказ по ID
func (r *InMemoryOrderRepository) GetOrderByID(orderID string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[orderID]
	if !ok {
		return nil, errors.New("order not found")
	}
	return order, nil
}
