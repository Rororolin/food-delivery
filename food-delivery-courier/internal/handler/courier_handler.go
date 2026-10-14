package handler

import (
	"encoding/json"
	"food-delivery-courier/internal/domain"
	"food-delivery-courier/internal/usecase"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CourierHandler struct {
	uc *usecase.CourierUseCase
}

func NewCourierHandler(uc *usecase.CourierUseCase) *CourierHandler {
	return &CourierHandler{uc: uc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /courier/orders/available
func (h *CourierHandler) GetAvailableOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.uc.GetAvailableOrders()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": orders})
}

// POST /courier/orders/{order_id}/respond
// Body: { "courier_id": "...", "action": "accept"|"reject", "reason": "..." }
func (h *CourierHandler) RespondToOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	var body struct {
		CourierID string `json:"courier_id"`
		Action    string `json:"action"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CourierID == "" {
		writeError(w, http.StatusBadRequest, "courier_id is required")
		return
	}

	switch body.Action {
	case "accept":
		order, err := h.uc.AcceptOrder(body.CourierID, orderID)
		if err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, order)
	case "reject":
		if err := h.uc.RejectOrder(body.CourierID, orderID, body.Reason); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
	default:
		writeError(w, http.StatusBadRequest, "action must be 'accept' or 'reject'")
	}
}

// GET /courier/orders/{order_id}/route?courier_id=...
func (h *CourierHandler) GetRoute(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")
	courierID := r.URL.Query().Get("courier_id")
	if courierID == "" {
		writeError(w, http.StatusBadRequest, "courier_id query param is required")
		return
	}

	route, err := h.uc.GetRoute(courierID, orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, route)
}

// POST /courier/orders/{order_id}/complete
// Body: { "courier_id": "..." }
func (h *CourierHandler) CompleteDelivery(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	var body struct {
		CourierID string `json:"courier_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CourierID == "" {
		writeError(w, http.StatusBadRequest, "courier_id is required")
		return
	}

	order, err := h.uc.CompleteDelivery(body.CourierID, orderID)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, order)
}

// GET /courier/orders/{order_id}
func (h *CourierHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")
	order, err := h.uc.GetOrder(orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, order)
}

// POST /courier/orders  — принять заказ от клиентского приложения
func (h *CourierHandler) ReceiveOrder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID                 string          `json:"id"`
		UserID             string          `json:"user_id"`
		CustomerName       string          `json:"customer_name"`
		RestaurantName     string          `json:"restaurant_name"`
		RestaurantAddress  string          `json:"restaurant_address"`
		RestaurantLocation domain.Location `json:"restaurant_location"`
		DeliveryAddress    string          `json:"delivery_address"`
		DeliveryLocation   domain.Location `json:"delivery_location"`
		ItemsCount         int             `json:"items_count"`
		DeliveryFee        float64         `json:"delivery_fee"`
		TotalAmount        float64         `json:"total_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id := body.ID
	if id == "" {
		id = uuid.NewString()
	}

	order := domain.Order{
		ID:                 id,
		UserID:             body.UserID,
		CustomerName:       body.CustomerName,
		RestaurantName:     body.RestaurantName,
		RestaurantAddress:  body.RestaurantAddress,
		RestaurantLocation: body.RestaurantLocation,
		DeliveryAddress:    body.DeliveryAddress,
		DeliveryLocation:   body.DeliveryLocation,
		ItemsCount:         body.ItemsCount,
		DeliveryFee:        body.DeliveryFee,
		TotalAmount:        body.TotalAmount,
		CreatedAt:          time.Now(),
	}

	if err := h.uc.AddOrderToPool(order); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"order_id": order.ID, "status": "available"})
}

// PATCH /courier/orders/{order_id}/sync — синхронизация статуса от админского сервиса
func (h *CourierHandler) SyncStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")
	var body struct {
		Status    string `json:"status"`
		CourierID string `json:"courier_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}

	switch body.Status {
	case "cancelled":
		if err := h.uc.ForceCancel(orderID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	case "assigned":
		if body.CourierID == "" {
			writeError(w, http.StatusBadRequest, "courier_id required for assigned status")
			return
		}
		_, err := h.uc.ForceAssign(orderID, body.CourierID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		_ = h.uc.SyncOrderStatus(orderID, domain.OrderStatus(body.Status))
	}

	writeJSON(w, http.StatusOK, map[string]string{"order_id": orderID, "status": body.Status})
}

// POST /courier/register
func (h *CourierHandler) RegisterCourier(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	courier := domain.Courier{
		ID:     uuid.NewString(),
		Name:   body.Name,
		Phone:  body.Phone,
		Status: domain.CourierStatusFree,
		CurrentLocation: domain.Location{Lat: 59.940, Lon: 30.350},
	}
	if err := h.uc.RegisterCourier(courier); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, courier)
}

// GET /courier/couriers  — список всех курьеров с полным ID (для админа и входа)
func (h *CourierHandler) GetAllCouriers(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.uc.GetAllCouriers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"couriers": couriers})
}

// GET /courier/find?name=...  — найти курьера по имени (для входа по имени)
func (h *CourierHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "name query param is required")
		return
	}
	courier, err := h.uc.FindCourierByName(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "Курьер с таким именем не найден")
		return
	}
	writeJSON(w, http.StatusOK, courier)
}

// GET /courier/profile?courier_id=...
func (h *CourierHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	courierID := r.URL.Query().Get("courier_id")
	if courierID == "" {
		writeError(w, http.StatusBadRequest, "courier_id is required")
		return
	}
	courier, err := h.uc.GetCourier(courierID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, courier)
}
