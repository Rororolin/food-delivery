package usecase

import (
	"errors"
	"food-delivery/internal/domain"
	"food-delivery/internal/port"
)

type catalogUseCase struct {
	repo port.CatalogRepository
}

// NewCatalogUseCase создаёт сценарий каталога с внедрённым репозиторием
func NewCatalogUseCase(repo port.CatalogRepository) port.CatalogUseCase {
	return &catalogUseCase{repo: repo}
}

// ListRestaurants возвращает рестораны, доставляющие по указанному адресу
func (uc *catalogUseCase) ListRestaurants(address string) ([]domain.Restaurant, error) {
	if address == "" {
		return nil, errors.New("address is required")
	}
	return uc.repo.GetRestaurantsByAddress(address)
}

// GetMenu возвращает меню ресторана
func (uc *catalogUseCase) GetMenu(restaurantID string) ([]domain.MenuItem, error) {
	if restaurantID == "" {
		return nil, errors.New("restaurantID is required")
	}
	return uc.repo.GetMenuByRestaurantID(restaurantID)
}
