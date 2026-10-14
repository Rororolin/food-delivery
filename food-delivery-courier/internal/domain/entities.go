package domain

import "time"

// ─── Статусы заказа (со стороны курьера) ────────────────────────────────────

type OrderStatus string

const (
	OrderStatusAvailable OrderStatus = "available"  // ждёт курьера
	OrderStatusAssigned  OrderStatus = "assigned"   // курьер принял
	OrderStatusCompleted OrderStatus = "completed"  // доставлен
	OrderStatusCancelled OrderStatus = "cancelled"  // отменён
)

// ─── Статус курьера ──────────────────────────────────────────────────────────

type CourierStatus string

const (
	CourierStatusFree CourierStatus = "free"
	CourierStatusBusy CourierStatus = "busy"
)

// ─── Сущности ────────────────────────────────────────────────────────────────

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Courier struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Phone           string        `json:"phone"`
	Status          CourierStatus `json:"status"`
	CurrentLocation Location      `json:"current_location"`
}

type Order struct {
	ID                 string      `json:"id"`
	Status             OrderStatus `json:"status"`
	UserID             string      `json:"user_id"`
	CustomerName       string      `json:"customer_name"`
	RestaurantName     string      `json:"restaurant_name"`
	RestaurantAddress  string      `json:"restaurant_address"`
	RestaurantLocation Location    `json:"restaurant_location"`
	DeliveryAddress    string      `json:"delivery_address"`
	DeliveryLocation   Location    `json:"delivery_location"`
	ItemsCount         int         `json:"items_count"`
	DeliveryFee        float64     `json:"delivery_fee"`
	TotalAmount        float64     `json:"total_amount"`
	AssignedCourierID  string      `json:"assigned_courier_id,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	CompletedAt        *time.Time  `json:"completed_at,omitempty"`
}

// ─── Маршрут ─────────────────────────────────────────────────────────────────

type Route struct {
	OrderID           string     `json:"order_id"`
	TotalDistanceKm   float64    `json:"total_distance_km"`
	TotalEstimatedMin int        `json:"total_estimated_minutes"`
	Legs              []RouteLeg `json:"legs"`
}

type RouteLeg struct {
	Type             string   `json:"type"` // "to_restaurant" / "to_client"
	From             Location `json:"from"`
	To               Location `json:"to"`
	DistanceKm       float64  `json:"distance_km"`
	EstimatedMinutes int      `json:"estimated_minutes"`
}
