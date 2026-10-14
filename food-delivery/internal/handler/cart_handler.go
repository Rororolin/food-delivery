package handler

import (
	"encoding/json"
	"food-delivery/internal/domain"
	"food-delivery/internal/port"
	"net/http"
)

// CartHandler — обработчик компонента корзины
type CartHandler struct {
	uc port.CartUseCase
}

func NewCartHandler(uc port.CartUseCase) *CartHandler {
	return &CartHandler{uc: uc}
}

type cartRequest struct {
	UserID   string             `json:"userId"`
	Action   domain.CartAction  `json:"action"`
	DishID   string             `json:"dishId"`
	Quantity int                `json:"quantity"`
}

// POST /cart
func (h *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	var req cartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	cart, err := h.uc.UpdateCart(req.UserID, req.Action, req.DishID, req.Quantity)
	if err != nil {
		switch err.Error() {
		case "dish not found: " + req.DishID:
			http.Error(w, `{"error":"dish not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		}
		return
	}

	writeJSON(w, http.StatusOK, cart)
}
