package payment

import (
	"fmt"
	"sync"
	"time"

	"food-delivery-admin/internal/domain"
)

// InMemoryPaymentGateway — имитация платёжного шлюза
type InMemoryPaymentGateway struct {
	mu      sync.RWMutex
	refunds map[string]*domain.Refund // ключ — order_id
}

func NewInMemoryPaymentGateway() *InMemoryPaymentGateway {
	return &InMemoryPaymentGateway{
		refunds: make(map[string]*domain.Refund),
	}
}

func (g *InMemoryPaymentGateway) ProcessRefund(orderID string, amount float64) (*domain.Refund, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.refunds[orderID]; exists {
		return nil, fmt.Errorf("refund for order %q already processed", orderID)
	}

	refund := &domain.Refund{
		OrderID:   orderID,
		Amount:    amount,
		Status:    domain.RefundStatusCompleted,
		CreatedAt: time.Now(),
	}
	g.refunds[orderID] = refund
	return refund, nil
}

func (g *InMemoryPaymentGateway) GetRefundStatus(orderID string) (*domain.Refund, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	r, ok := g.refunds[orderID]
	if !ok {
		return nil, fmt.Errorf("refund for order %q not found", orderID)
	}
	cp := *r
	return &cp, nil
}
