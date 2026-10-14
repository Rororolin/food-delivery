package repository

import (
	"errors"
	"food-delivery/internal/domain"
	"sync"
)

// InMemoryCartRepository — хранилище корзин в памяти
type InMemoryCartRepository struct {
	mu    sync.RWMutex
	carts map[string]*domain.Cart
}

func NewInMemoryCartRepository() *InMemoryCartRepository {
	return &InMemoryCartRepository{
		carts: make(map[string]*domain.Cart),
	}
}

// GetCart возвращает корзину пользователя. Если её нет — возвращает пустую.
func (r *InMemoryCartRepository) GetCart(userID string) (*domain.Cart, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart, ok := r.carts[userID]
	if !ok {
		return &domain.Cart{UserID: userID, Items: []domain.CartItem{}}, nil
	}
	return cart, nil
}

// SaveCart сохраняет корзину пользователя
func (r *InMemoryCartRepository) SaveCart(cart *domain.Cart) error {
	if cart == nil {
		return errors.New("cart must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.carts[cart.UserID] = cart
	return nil
}
