package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"food-delivery-courier/internal/domain"
	"food-delivery-courier/internal/port"
	"net/http"
	"time"
)

const adminAPI = "http://localhost:8081"

type CourierUseCase struct {
	distRepo    port.OrderDistributionRepository
	orderRepo   port.OrderRepository
	courierRepo port.CourierStatsRepository
	notifySvc   port.NotificationService
	geoSvc      port.GeoService
}

func NewCourierUseCase(
	distRepo port.OrderDistributionRepository,
	orderRepo port.OrderRepository,
	courierRepo port.CourierStatsRepository,
	notifySvc port.NotificationService,
	geoSvc port.GeoService,
) *CourierUseCase {
	return &CourierUseCase{
		distRepo:    distRepo,
		orderRepo:   orderRepo,
		courierRepo: courierRepo,
		notifySvc:   notifySvc,
		geoSvc:      geoSvc,
	}
}

// syncAdmin отправляет обновление статуса в админский сервис (не блокирует при ошибке)
func syncAdmin(orderID, status, courierID string) {
	body, _ := json.Marshal(map[string]string{
		"status":     status,
		"courier_id": courierID,
	})
	req, err := http.NewRequest("PATCH", adminAPI+"/admin/orders/"+orderID+"/sync", bytes.NewBuffer(body))
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

// ─── 1. Получить доступные заказы ────────────────────────────────────────────
// GET /courier/orders/available
func (uc *CourierUseCase) GetAvailableOrders() ([]domain.Order, error) {
	return uc.distRepo.GetAvailableOrders()
}

// ─── 2. Принять заказ ────────────────────────────────────────────────────────
// POST /courier/orders/{order_id}/respond  { action: "accept" }
func (uc *CourierUseCase) AcceptOrder(courierID, orderID string) (domain.Order, error) {
	courier, err := uc.courierRepo.GetCourier(courierID)
	if err != nil {
		return domain.Order{}, errors.New("courier not found")
	}
	if courier.Status == domain.CourierStatusBusy {
		return domain.Order{}, errors.New("courier is already busy")
	}

	// Убираем из пула доступных
	if err := uc.distRepo.RemoveAvailableOrder(orderID); err != nil {
		return domain.Order{}, errors.New("order is no longer available")
	}

	// Назначаем курьера
	if err := uc.orderRepo.AssignCourier(orderID, courierID); err != nil {
		return domain.Order{}, err
	}

	// Обновляем статус курьера
	courier.Status = domain.CourierStatusBusy
	_ = uc.courierRepo.SaveCourier(courier)

	_ = uc.notifySvc.SendOrderAccepted(courierID, orderID)

	// Синхронизируем статус с админским сервисом
	go syncAdmin(orderID, "assigned", courierID)

	return uc.orderRepo.FindByID(orderID)
}

// ─── 3. Отказаться от заказа ─────────────────────────────────────────────────
// POST /courier/orders/{order_id}/respond  { action: "reject", reason: "..." }
func (uc *CourierUseCase) RejectOrder(courierID, orderID, reason string) error {
	_ = uc.notifySvc.SendOrderRejected(courierID, orderID, reason)
	return nil
}

// ─── 4. Построить маршрут ────────────────────────────────────────────────────
// GET /courier/orders/{order_id}/route
func (uc *CourierUseCase) GetRoute(courierID, orderID string) (domain.Route, error) {
	courier, err := uc.courierRepo.GetCourier(courierID)
	if err != nil {
		return domain.Route{}, errors.New("courier not found")
	}

	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return domain.Route{}, err
	}

	route, err := uc.geoSvc.GetDeliveryRoute(
		courier.CurrentLocation,
		order.RestaurantLocation,
		order.DeliveryLocation,
	)
	if err != nil {
		return domain.Route{}, err
	}
	route.OrderID = orderID
	return route, nil
}

// ─── 5. Подтвердить доставку ─────────────────────────────────────────────────
// POST /courier/orders/{order_id}/complete
func (uc *CourierUseCase) CompleteDelivery(courierID, orderID string) (domain.Order, error) {
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if order.AssignedCourierID != courierID {
		return domain.Order{}, errors.New("this order is not assigned to you")
	}
	if order.Status != domain.OrderStatusAssigned {
		return domain.Order{}, errors.New("order is not in assigned status")
	}

	if err := uc.orderRepo.UpdateStatus(orderID, domain.OrderStatusCompleted); err != nil {
		return domain.Order{}, err
	}

	// Освобождаем курьера
	courier, _ := uc.courierRepo.GetCourier(courierID)
	courier.Status = domain.CourierStatusFree
	_ = uc.courierRepo.SaveCourier(courier)

	_ = uc.notifySvc.SendOrderCompleted(order.UserID, orderID)

	// Синхронизируем статус с админским сервисом
	go syncAdmin(orderID, "completed", courierID)

	now := time.Now()
	order.Status = domain.OrderStatusCompleted
	order.CompletedAt = &now
	return order, nil
}

func (uc *CourierUseCase) GetAllCouriers() ([]domain.Courier, error) {
	return uc.courierRepo.GetAllCouriers()
}

// ─── Добавить заказ в пул (вызывается клиентским сервером) ───────────────────
func (uc *CourierUseCase) AddOrderToPool(order domain.Order) error {
	order.Status = domain.OrderStatusAvailable
	if err := uc.orderRepo.Save(order); err != nil {
		return err
	}
	_ = uc.notifySvc.SendNewOrderAvailable(order.ID)
	return uc.distRepo.AddAvailableOrder(order)
}

// ─── Получить заказ по ID ────────────────────────────────────────────────────
func (uc *CourierUseCase) GetOrder(orderID string) (domain.Order, error) {
	return uc.orderRepo.FindByID(orderID)
}

// ─── Регистрация / получение курьера ─────────────────────────────────────────
func (uc *CourierUseCase) RegisterCourier(courier domain.Courier) error {
	return uc.courierRepo.SaveCourier(courier)
}

func (uc *CourierUseCase) GetCourier(courierID string) (domain.Courier, error) {
	return uc.courierRepo.GetCourier(courierID)
}

// ForceAssign — принудительно назначает курьера на заказ (вызывается от админа)
func (uc *CourierUseCase) ForceAssign(orderID, courierID string) (domain.Order, error) {
	courier, err := uc.courierRepo.GetCourier(courierID)
	if err != nil {
		return domain.Order{}, errors.New("courier not found")
	}

	// Убираем из пула доступных (если ещё там)
	_ = uc.distRepo.RemoveAvailableOrder(orderID)

	// Назначаем курьера
	if err := uc.orderRepo.AssignCourier(orderID, courierID); err != nil {
		return domain.Order{}, err
	}

	// Помечаем курьера занятым
	courier.Status = domain.CourierStatusBusy
	_ = uc.courierRepo.SaveCourier(courier)

	return uc.orderRepo.FindByID(orderID)
}

// ForceCancel — принудительно отменяет заказ (вызывается от админа), освобождает курьера
func (uc *CourierUseCase) ForceCancel(orderID string) error {
	order, err := uc.orderRepo.FindByID(orderID)
	if err != nil {
		// Заказа нет на курьерском сервере — ок
		return nil
	}

	// Освобождаем курьера если был назначен
	if order.AssignedCourierID != "" {
		courier, err := uc.courierRepo.GetCourier(order.AssignedCourierID)
		if err == nil {
			courier.Status = domain.CourierStatusFree
			_ = uc.courierRepo.SaveCourier(courier)
		}
	}

	// Убираем из пула доступных
	_ = uc.distRepo.RemoveAvailableOrder(orderID)

	return uc.orderRepo.UpdateStatus(orderID, domain.OrderStatusCancelled)
}

// SyncOrderStatus — обновляет статус заказа от внешнего сервиса (без проверок)
func (uc *CourierUseCase) SyncOrderStatus(orderID string, status domain.OrderStatus) error {
	return uc.orderRepo.UpdateStatus(orderID, status)
}

// FindCourierByName — поиск курьера по имени (для входа по имени)
func (uc *CourierUseCase) FindCourierByName(name string) (domain.Courier, error) {
	return uc.courierRepo.FindByName(name)
}
