package repository

import (
	"errors"
	"food-delivery/internal/domain"
	"sync"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]domain.User)}
}

func (r *InMemoryUserRepository) Save(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetAll() ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *InMemoryUserRepository) GetByID(id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	return u, nil
}
