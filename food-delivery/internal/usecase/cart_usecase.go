package usecase

import (
	"errors"
	"food-delivery/internal/domain"
	"food-delivery/internal/port"
)

type cartUseCase struct {
	cartRepo    port.CartRepository
	catalogRepo port.CatalogRepository
}

// NewCartUseCase создаёт сценарий корзины
func NewCartUseCase(cartRepo port.CartRepository, catalogRepo port.CatalogRepository) port.CartUseCase {
	return &cartUseCase{cartRepo: cartRepo, catalogRepo: catalogRepo}
}

// GetCart возвращает текущую корзину пользователя
func (uc *cartUseCase) GetCart(userID string) (*domain.Cart, error) {
	return uc.cartRepo.GetCart(userID)
}

// UpdateCart добавляет или удаляет блюдо из корзины
func (uc *cartUseCase) UpdateCart(userID string, action domain.CartAction, dishID string, quantity int) (*domain.Cart, error) {
	if userID == "" || dishID == "" {
		return nil, errors.New("userID and dishID are required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	cart, err := uc.cartRepo.GetCart(userID)
	if err != nil {
		return nil, err
	}

	switch action {
	case domain.CartActionAdd:
		if err := uc.addItem(cart, dishID, quantity); err != nil {
			return nil, err
		}
	case domain.CartActionRemove:
		uc.removeItem(cart, dishID, quantity)
	default:
		return nil, errors.New("unknown action: use 'add' or 'remove'")
	}

	cart.TotalPrice = calcTotal(cart.Items)

	if err := uc.cartRepo.SaveCart(cart); err != nil {
		return nil, err
	}
	return cart, nil
}

// addItem находит блюдо в каталоге и добавляет его в корзину
func (uc *cartUseCase) addItem(cart *domain.Cart, dishID string, quantity int) error {
	// Ищем блюдо во всех меню (простая реализация для in-memory)
	dish, err := uc.findDish(dishID)
	if err != nil {
		return err
	}

	for i, item := range cart.Items {
		if item.DishID == dishID {
			cart.Items[i].Quantity += quantity
			return nil
		}
	}
	cart.Items = append(cart.Items, domain.CartItem{
		DishID:   dish.ID,
		Name:     dish.Name,
		Price:    dish.Price,
		Quantity: quantity,
	})
	return nil
}

// removeItem уменьшает количество блюда или убирает его из корзины
func (uc *cartUseCase) removeItem(cart *domain.Cart, dishID string, quantity int) {
	for i, item := range cart.Items {
		if item.DishID == dishID {
			cart.Items[i].Quantity -= quantity
			if cart.Items[i].Quantity <= 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			}
			return
		}
	}
}

// findDish перебирает меню ресторанов в поисках блюда по ID
func (uc *cartUseCase) findDish(dishID string) (*domain.MenuItem, error) {
	// В реальном проекте каталог передавал бы метод поиска по dishID напрямую.
	// Здесь используем GetMenuByRestaurantID через известные ID ресторанов.
	for _, restaurantID := range []string{"r1", "r2", "r3", "r4", "r5", "r6", "r7", "r8"} {
		menu, err := uc.catalogRepo.GetMenuByRestaurantID(restaurantID)
		if err != nil {
			continue
		}
		for _, item := range menu {
			if item.ID == dishID {
				return &item, nil
			}
		}
	}
	return nil, errors.New("dish not found: " + dishID)
}

func calcTotal(items []domain.CartItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
