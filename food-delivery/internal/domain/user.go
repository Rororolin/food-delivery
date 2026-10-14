package domain

import "time"

// User — клиент сервиса
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"createdAt"`
}
