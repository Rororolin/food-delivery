package repository

import (
	"errors"
	"food-delivery-courier/internal/domain"
	"sync"
)

// ─── Репозиторий заказов ──────────────────────────────────────────────────────

type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]domain.Order),
	}
}

func (r *InMemoryOrderRepository) Save(order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) FindByID(orderID string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[orderID]
	if !ok {
		return domain.Order{}, errors.New("order not found: " + orderID)
	}
	return o, nil
}

func (r *InMemoryOrderRepository) UpdateStatus(orderID string, status domain.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[orderID]
	if !ok {
		return errors.New("order not found: " + orderID)
	}
	o.Status = status
	r.orders[orderID] = o
	return nil
}

func (r *InMemoryOrderRepository) AssignCourier(orderID, courierID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[orderID]
	if !ok {
		return errors.New("order not found: " + orderID)
	}
	o.AssignedCourierID = courierID
	o.Status = domain.OrderStatusAssigned
	r.orders[orderID] = o
	return nil
}

func (r *InMemoryOrderRepository) GetAll() ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Order, 0, len(r.orders))
	for _, o := range r.orders {
		result = append(result, o)
	}
	return result, nil
}

// ─── Репозиторий распределения (кэш доступных заказов) ───────────────────────

type InMemoryOrderDistributionRepository struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewInMemoryOrderDistributionRepository() *InMemoryOrderDistributionRepository {
	return &InMemoryOrderDistributionRepository{
		orders: make(map[string]domain.Order),
	}
}

func (r *InMemoryOrderDistributionRepository) GetAvailableOrders() ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Order, 0, len(r.orders))
	for _, o := range r.orders {
		result = append(result, o)
	}
	return result, nil
}

func (r *InMemoryOrderDistributionRepository) AddAvailableOrder(order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderDistributionRepository) RemoveAvailableOrder(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.orders, orderID)
	return nil
}
