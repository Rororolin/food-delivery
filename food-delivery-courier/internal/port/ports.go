package port

import "food-delivery-courier/internal/domain"

// ─── Репозиторий распределения заказов (кэш доступных заказов) ───────────────
type OrderDistributionRepository interface {
	GetAvailableOrders() ([]domain.Order, error)
	AddAvailableOrder(order domain.Order) error
	RemoveAvailableOrder(orderID string) error
}

// ─── Репозиторий заказов ──────────────────────────────────────────────────────
type OrderRepository interface {
	FindByID(orderID string) (domain.Order, error)
	UpdateStatus(orderID string, status domain.OrderStatus) error
	AssignCourier(orderID, courierID string) error
	Save(order domain.Order) error
	GetAll() ([]domain.Order, error)
}

// ─── Репозиторий статистики курьера ──────────────────────────────────────────
type CourierStatsRepository interface {
	GetCourier(courierID string) (domain.Courier, error)
	FindByName(name string) (domain.Courier, error)
	SaveCourier(courier domain.Courier) error
	GetAllCouriers() ([]domain.Courier, error)
}

// ─── Сервис уведомлений ───────────────────────────────────────────────────────
type NotificationService interface {
	SendOrderAccepted(courierID, orderID string) error
	SendOrderRejected(courierID, orderID, reason string) error
	SendOrderCompleted(userID, orderID string) error
	SendNewOrderAvailable(orderID string) error
}

// ─── Гео-сервис (построение маршрута) ────────────────────────────────────────
type GeoService interface {
	GetDeliveryRoute(courier, restaurant, delivery domain.Location) (domain.Route, error)
	CalculateDistance(from, to domain.Location) float64
}
