package notification

import (
	"fmt"
	"log"
)

// LogNotificationService — реализация NotificationService через логирование
// (имитирует отправку уведомлений без реального транспорта)
type LogNotificationService struct{}

func NewLogNotificationService() *LogNotificationService {
	return &LogNotificationService{}
}

func (s *LogNotificationService) NotifyCustomer(customerID, message string) error {
	if customerID == "" {
		return fmt.Errorf("customerID is empty")
	}
	log.Printf("[NOTIFICATION] -> Customer %s: %s", customerID, message)
	return nil
}

func (s *LogNotificationService) NotifyCourier(courierID, message string) error {
	if courierID == "" {
		// Курьер может быть не назначен — просто пропускаем
		return nil
	}
	log.Printf("[NOTIFICATION] -> Courier %s: %s", courierID, message)
	return nil
}
