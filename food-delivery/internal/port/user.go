package port

import "food-delivery/internal/domain"

type UserRepository interface {
	Save(user domain.User) error
	GetAll() ([]domain.User, error)
	GetByID(id string) (domain.User, error)
}
