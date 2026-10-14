package usecase

import (
	"errors"
	"food-delivery/internal/domain"
	"food-delivery/internal/port"
	"time"

	"github.com/google/uuid"
)

type orderUseCase struct {
	orderRepo port.OrderRepository
	userRepo  port.UserRepository
}

// NewOrderUseCase создаёт сценарий оформления и отслеживания заказа
func NewOrderUseCase(orderRepo port.OrderRepository, userRepo port.UserRepository) port.OrderUseCase {
	return &orderUseCase{orderRepo: orderRepo, userRepo: userRepo}
}

// PlaceOrder создаёт новый заказ
func (uc *orderUseCase) PlaceOrder(userID string, items []domain.OrderItem, address, paymentMethod string) (*domain.Order, error) {
	if userID == "" || address == "" || paymentMethod == "" {
		return nil, errors.New("userID, address and paymentMethod are required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	// Проверяем, не заблокирован ли пользователь
	user, err := uc.userRepo.GetByID(userID)
	if err == nil && user.IsBlocked {
		return nil, errors.New("ваш аккаунт заблокирован. Оформление заказа недоступно")
	}

	now := time.Now()
	order := &domain.Order{
		ID:                uuid.NewString(),
		UserID:            userID,
		Items:             items,
		Address:           address,
		PaymentMethod:     paymentMethod,
		Status:            domain.OrderStatusCreated,
		EstimatedDelivery: now.Add(45 * time.Minute),
		CreatedAt:         now,
	}

	if err := uc.orderRepo.SaveOrder(order); err != nil {
		return nil, err
	}
	return order, nil
}

// TrackOrder возвращает текущий статус заказа
func (uc *orderUseCase) TrackOrder(orderID string) (*domain.Order, error) {
	if orderID == "" {
		return nil, errors.New("orderID is required")
	}
	return uc.orderRepo.GetOrderByID(orderID)
}
