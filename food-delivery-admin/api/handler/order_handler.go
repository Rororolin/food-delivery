package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"food-delivery-admin/internal/domain"
	"food-delivery-admin/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const courierAPI = "http://localhost:8082"

// syncCourier отправляет обновление статуса на курьерский сервис
func syncCourier(orderID, status, courierID string) {
	payload := map[string]string{"status": status}
	if courierID != "" {
		payload["courier_id"] = courierID
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PATCH",
		courierAPI+"/courier/orders/"+orderID+"/sync",
		bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}

// OrderHandler — REST-обработчики для Сервиса управления заказами
type OrderHandler struct {
	uc *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

// SyncStatus PATCH /admin/orders/{order_id}/sync — синхронизация статуса от курьерского сервиса
func (h *OrderHandler) SyncStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	var body struct {
		Status    string `json:"status"`
		CourierID string `json:"courier_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}

	// Маппинг статусов курьерского сервиса → статусы админского
	statusMap := map[string]domain.OrderStatus{
		"assigned":  domain.OrderStatusAssigned,
		"completed": domain.OrderStatusDelivered,
		"cancelled": domain.OrderStatusCancelled,
	}

	newStatus, ok := statusMap[body.Status]
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown status: "+body.Status)
		return
	}

	// Если нужно — назначаем курьера
	if body.CourierID != "" {
		_, _ = h.uc.ReassignCourier(orderID, body.CourierID)
	}

	order, err := h.uc.ChangeOrderStatus(orderID, newStatus)
	if err != nil {
		// Если заказ уже в нужном статусе — не ошибка
		writeJSON(w, http.StatusOK, map[string]string{"order_id": orderID, "status": string(newStatus), "note": "already synced"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"order_id": order.ID, "new_status": string(order.Status)})
}

// CreateOrder POST /admin/orders — приём заказа от клиентского приложения
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID          string  `json:"id"`
		CustomerID  string  `json:"customer_id"`
		TotalAmount float64 `json:"totalAmount"`
		IsPaid      bool    `json:"is_paid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, "fields id and customer_id are required")
		return
	}

	order, err := h.uc.CreateOrder(body.ID, body.CustomerID, body.TotalAmount, body.IsPaid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

// GetOrders GET /admin/orders
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	orders, err := h.uc.GetOrders(status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"orders": orders})
}

// GetOrder GET /admin/orders/{order_id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	order, err := h.uc.GetOrderByID(orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// UpdateOrderStatus PATCH /admin/orders/{order_id}/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, http.StatusBadRequest, "field 'status' is required")
		return
	}

	newStatus := domain.OrderStatus(body.Status)
	order, err := h.uc.ChangeOrderStatus(orderID, newStatus)
	if err != nil {
		// Различаем «не найден» и «недопустимый переход»
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"order_id":   order.ID,
		"new_status": order.Status,
	})
}

// CancelOrder POST /admin/orders/{order_id}/cancel
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	order, err := h.uc.CancelOrder(orderID)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	// Синхронизируем отмену с курьерским сервисом
	go syncCourier(orderID, "cancelled", "")

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"order_id":   order.ID,
		"new_status": order.Status,
	})
}

// ReassignCourier PATCH /admin/orders/{order_id}/courier
func (h *OrderHandler) ReassignCourier(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	var body struct {
		CourierID string `json:"courier_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CourierID == "" {
		writeError(w, http.StatusBadRequest, "field 'courier_id' is required")
		return
	}

	order, err := h.uc.ReassignCourier(orderID, body.CourierID)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	// Синхронизируем назначение с курьерским сервисом — заказ появится у курьера в активных
	go syncCourier(orderID, "assigned", body.CourierID)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"order_id":   order.ID,
		"courier_id": order.CourierID,
	})
}

// ProcessRefund POST /admin/orders/{order_id}/refund
func (h *OrderHandler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "order_id")

	refund, err := h.uc.ProcessRefund(orderID)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, refund)
}
