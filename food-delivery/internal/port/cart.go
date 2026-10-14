package port

import "food-delivery/internal/domain"

// CartRepository — порт для хранилища корзин
type CartRepository interface {
	GetCart(userID string) (*domain.Cart, error)
	SaveCart(cart *domain.Cart) error
}

// CartUseCase — порт для сценариев корзины
type CartUseCase interface {
	UpdateCart(userID string, action domain.CartAction, dishID string, quantity int) (*domain.Cart, error)
	GetCart(userID string) (*domain.Cart, error)
}
