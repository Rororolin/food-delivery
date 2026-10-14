package handler

import (
	"encoding/json"
	"food-delivery/internal/domain"
	"food-delivery/internal/port"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// OrderHandler — обработчик компонентов оформления и отслеживания заказа
type OrderHandler struct {
	uc port.OrderUseCase
}

func NewOrderHandler(uc port.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

type placeOrderRequest struct {
	UserID        string               `json:"userId"`
	Cart          []domain.OrderItem   `json:"cart"`
	Address       string               `json:"address"`
	PaymentMethod string               `json:"paymentMethod"`
}

// POST /orders — оформление заказа
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req placeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	order, err := h.uc.PlaceOrder(req.UserID, req.Cart, req.Address, req.PaymentMethod)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"orderId":           order.ID,
		"status":            order.Status,
		"estimatedDelivery": order.EstimatedDelivery,
	})
}

// GET /orders/{order_id} — отслеживание заказа
func (h *OrderHandler) TrackOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	order, err := h.uc.TrackOrder(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			http.Error(w, `{"error":"order not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"orderId":           order.ID,
		"status":            order.Status,
		"estimatedDelivery": order.EstimatedDelivery,
		"createdAt":         order.CreatedAt,
	})
}
