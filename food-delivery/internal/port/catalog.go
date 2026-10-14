package port

import "food-delivery/internal/domain"

// CatalogRepository — порт для хранилища каталога
type CatalogRepository interface {
	GetRestaurantsByAddress(address string) ([]domain.Restaurant, error)
	GetMenuByRestaurantID(restaurantID string) ([]domain.MenuItem, error)
}

// CatalogUseCase — порт для сценариев каталога
type CatalogUseCase interface {
	ListRestaurants(address string) ([]domain.Restaurant, error)
	GetMenu(restaurantID string) ([]domain.MenuItem, error)
}
