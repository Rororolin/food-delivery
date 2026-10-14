package domain

// Restaurant — ресторан в каталоге
type Restaurant struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Cuisine      string  `json:"cuisine"`
	Rating       float64 `json:"rating"`
	DeliveryTime int     `json:"deliveryTime"` // минуты
}

// MenuItem — блюдо в меню ресторана
type MenuItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
