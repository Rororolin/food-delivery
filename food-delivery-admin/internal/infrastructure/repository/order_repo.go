package repository

import (
	"fmt"
	"sync"

	"food-delivery-admin/internal/domain"

	"github.com/google/uuid"
)

// InMemoryOrderRepository — in-memory реализация OrderRepository
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *InMemoryOrderRepository) GetAll(status string) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Order
	for _, o := range r.orders {
		if status == "" || string(o.Status) == status {
			cp := *o
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r *InMemoryOrderRepository) GetByID(id string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %q not found", id)
	}
	cp := *o
	return &cp, nil
}

func (r *InMemoryOrderRepository) UpdateStatus(id string, newStatus domain.OrderStatus) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %q not found", id)
	}
	o.Status = newStatus
	cp := *o
	return &cp, nil
}

func (r *InMemoryOrderRepository) AssignCourier(orderID, courierID string) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, ok := r.orders[orderID]
	if !ok {
		return nil, fmt.Errorf("order %q not found", orderID)
	}
	o.CourierID = courierID
	cp := *o
	return &cp, nil
}

func (r *InMemoryOrderRepository) Save(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.ID == "" {
		order.ID = uuid.NewString()
	}
	cp := *order
	r.orders[order.ID] = &cp
	return nil
}

func (r *InMemoryOrderRepository) MarkRefunded(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	o, ok := r.orders[orderID]
	if !ok {
		return fmt.Errorf("order %q not found", orderID)
	}
	o.IsRefunded = true
	return nil
}
