package repository

import (
	"errors"
	"food-delivery-courier/internal/domain"
	"strings"
	"sync"
)

type InMemoryCourierStatsRepository struct {
	mu       sync.RWMutex
	couriers map[string]domain.Courier
}

func NewInMemoryCourierStatsRepository() *InMemoryCourierStatsRepository {
	return &InMemoryCourierStatsRepository{
		couriers: make(map[string]domain.Courier),
	}
}

func (r *InMemoryCourierStatsRepository) GetCourier(courierID string) (domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.couriers[courierID]
	if !ok {
		return domain.Courier{}, errors.New("courier not found: " + courierID)
	}
	return c, nil
}

// FindByName ищет курьера по имени (регистронезависимо)
func (r *InMemoryCourierStatsRepository) FindByName(name string) (domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	nameLower := strings.ToLower(strings.TrimSpace(name))
	for _, c := range r.couriers {
		if strings.ToLower(c.Name) == nameLower {
			return c, nil
		}
	}
	return domain.Courier{}, errors.New("courier not found by name: " + name)
}

func (r *InMemoryCourierStatsRepository) SaveCourier(courier domain.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.couriers[courier.ID] = courier
	return nil
}

func (r *InMemoryCourierStatsRepository) GetAllCouriers() ([]domain.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Courier, 0, len(r.couriers))
	for _, c := range r.couriers {
		result = append(result, c)
	}
	return result, nil
}
