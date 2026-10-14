package repository

import (
	"fmt"
	"sync"

	"food-delivery-admin/internal/domain"

	"github.com/google/uuid"
)

// InMemoryUserRepository — in-memory реализация UserRepository
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) GetAll() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.User
	for _, u := range r.users {
		cp := *u
		result = append(result, &cp)
	}
	return result, nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %q not found", id)
	}
	cp := *u
	return &cp, nil
}

func (r *InMemoryUserRepository) Block(id string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %q not found", id)
	}
	u.IsBlocked = true
	cp := *u
	return &cp, nil
}

func (r *InMemoryUserRepository) Unblock(id string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user %q not found", id)
	}
	u.IsBlocked = false
	cp := *u
	return &cp, nil
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	cp := *user
	r.users[user.ID] = &cp
	return nil
}
