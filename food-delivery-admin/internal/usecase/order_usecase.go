package usecase

import (
	"fmt"
	"time"

	"food-delivery-admin/internal/domain"
	"food-delivery-admin/internal/ports"
)

// OrderUseCase — сценарии приложения для управления заказами
type OrderUseCase struct {
	orderRepo   ports.OrderRepository
	paymentGW   ports.PaymentGateway
	notifySvc   ports.NotificationService
}

func NewOrderUseCase(
	orderRepo ports.OrderRepository,
	paymentGW ports.PaymentGateway,
	notifySvc ports.NotificationService,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
		paymentGW: paymentGW,
		notifySvc: notifySvc,
	}
}

// CreateOrder — создание заказа (вызывается клиентским приложением)
func (uc *OrderUseCase) CreateOrder(id, customerID string, totalAmount float64, isPaid bool) (*domain.Order, error) {
	order := &domain.Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      domain.OrderStatusPending,
		TotalAmount: totalAmount,
		IsPaid:      isPaid,
		CreatedAt:   time.Now(),
	}
	if err := uc.orderRepo.Save(order); err != nil {
		return nil, err
	}
	return order, nil
}

// GetOrders — получение списка заказов с необязательным фильтром по статусу
func (uc *OrderUseCase) GetOrders(status string) ([]*domain.Order, error) {
	return uc.orderRepo.GetAll(status)
}

// GetOrderByID — получение заказа по идентификатору
func (uc *OrderUseCase) GetOrderByID(orderID string) (*domain.Order, error) {
	return uc.orderRepo.GetByID(orderID)
}

// ChangeOrderStatus — изменение статуса заказа администратором
func (uc *OrderUseCase) ChangeOrderStatus(orderID string, newStatus domain.OrderStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if !order.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("invalid transition from %q to %q", order.Status, newStatus)
	}

	updated, err := uc.orderRepo.UpdateStatus(orderID, newStatus)
	if err != nil {
		return nil, err
	}

	// Уведомляем участников об изменении статуса
	msg := fmt.Sprintf("Статус заказа %s изменён на: %s", orderID, newStatus)
	_ = uc.notifySvc.NotifyCustomer(updated.CustomerID, msg)
	_ = uc.notifySvc.NotifyCourier(updated.CourierID, msg)

	return updated, nil
}

// CancelOrder — отмена заказа с возвратом средств (если заказ оплачен)
func (uc *OrderUseCase) CancelOrder(orderID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if !order.CanTransitionTo(domain.OrderStatusCancelled) {
		return nil, fmt.Errorf("cannot cancel order in status %q", order.Status)
	}

	updated, err := uc.orderRepo.UpdateStatus(orderID, domain.OrderStatusCancelled)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Заказ %s отменён администратором", orderID)
	_ = uc.notifySvc.NotifyCustomer(updated.CustomerID, msg)
	_ = uc.notifySvc.NotifyCourier(updated.CourierID, msg)

	return updated, nil
}

// ReassignCourier — переназначение курьера на заказ
func (uc *OrderUseCase) ReassignCourier(orderID, newCourierID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.Status == domain.OrderStatusDelivered || order.Status == domain.OrderStatusCancelled {
		return nil, fmt.Errorf("cannot reassign courier for order in status %q", order.Status)
	}

	oldCourierID := order.CourierID
	updated, err := uc.orderRepo.AssignCourier(orderID, newCourierID)
	if err != nil {
		return nil, err
	}

	_ = uc.notifySvc.NotifyCustomer(updated.CustomerID,
		fmt.Sprintf("Курьер по заказу %s переназначен", orderID))
	if oldCourierID != "" {
		_ = uc.notifySvc.NotifyCourier(oldCourierID,
			fmt.Sprintf("Вы сняты с заказа %s", orderID))
	}
	_ = uc.notifySvc.NotifyCourier(newCourierID,
		fmt.Sprintf("Вам назначен заказ %s", orderID))

	return updated, nil
}

// ProcessRefund — явный возврат средств по заказу
func (uc *OrderUseCase) ProcessRefund(orderID string) (*domain.Refund, error) {
	order, err := uc.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.IsRefunded {
		return nil, fmt.Errorf("refund for order %q already processed", orderID)
	}

	refund, err := uc.paymentGW.ProcessRefund(orderID, order.TotalAmount)
	if err != nil {
		return nil, err
	}

	_ = uc.orderRepo.MarkRefunded(orderID)

	_ = uc.notifySvc.NotifyCustomer(order.CustomerID,
		fmt.Sprintf("Возврат %.2f руб. по заказу %s выполнен", order.TotalAmount, orderID))

	return refund, nil
}
