package domain

// CartItem — позиция в корзине
type CartItem struct {
	DishID   string  `json:"dishId"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Cart — корзина пользователя
type Cart struct {
	UserID     string     `json:"userId"`
	Items      []CartItem `json:"cart"`
	TotalPrice float64    `json:"totalPrice"`
}

// CartAction — действие с корзиной
type CartAction string

const (
	CartActionAdd    CartAction = "add"
	CartActionRemove CartAction = "remove"
)
