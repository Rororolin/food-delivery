package usecase

import (
	"fmt"
	"time"
	"food-delivery-admin/internal/domain"
	"food-delivery-admin/internal/ports"
)

// UserUseCase — сценарии приложения для управления пользователями
type UserUseCase struct {
	userRepo  ports.UserRepository
	notifySvc ports.NotificationService
}

func NewUserUseCase(userRepo ports.UserRepository, notifySvc ports.NotificationService) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, notifySvc: notifySvc}
}

// CreateUser — создание пользователя (вызывается при синхронизации с клиентским сервисом)
func (uc *UserUseCase) CreateUser(id, name, phone string) (*domain.User, error) {
	user := &domain.User{
		ID:        id,
		Name:      name,
		Email:     phone, // phone хранится в поле Email как контакт
		Role:      domain.RoleCustomer,
		IsBlocked: false,
		CreatedAt: time.Now(),
	}
	if err := uc.userRepo.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers — получение списка всех пользователей
func (uc *UserUseCase) GetUsers() ([]*domain.User, error) {
	return uc.userRepo.GetAll()
}

// GetUserByID — получение пользователя по идентификатору
func (uc *UserUseCase) GetUserByID(userID string) (*domain.User, error) {
	return uc.userRepo.GetByID(userID)
}

// BlockUser — блокировка пользователя
func (uc *UserUseCase) BlockUser(userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user.IsBlocked {
		return nil, fmt.Errorf("user %q is already blocked", userID)
	}

	updated, err := uc.userRepo.Block(userID)
	if err != nil {
		return nil, err
	}

	_ = uc.notifySvc.NotifyCustomer(userID, "Ваш аккаунт заблокирован администратором")
	return updated, nil
}

// UnblockUser — разблокировка пользователя
func (uc *UserUseCase) UnblockUser(userID string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if !user.IsBlocked {
		return nil, fmt.Errorf("user %q is not blocked", userID)
	}

	updated, err := uc.userRepo.Unblock(userID)
	if err != nil {
		return nil, err
	}

	_ = uc.notifySvc.NotifyCustomer(userID, "Ваш аккаунт разблокирован администратором")
	return updated, nil
}
