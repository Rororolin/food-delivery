package notification

import "log"

type MockNotificationService struct{}

func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{}
}

func (s *MockNotificationService) SendOrderAccepted(courierID, orderID string) error {
	log.Printf("[NOTIFY] Курьер %s принял заказ %s", courierID, orderID)
	return nil
}

func (s *MockNotificationService) SendOrderRejected(courierID, orderID, reason string) error {
	log.Printf("[NOTIFY] Курьер %s отказался от заказа %s: %s", courierID, orderID, reason)
	return nil
}

func (s *MockNotificationService) SendOrderCompleted(userID, orderID string) error {
	log.Printf("[NOTIFY] Пользователь %s: заказ %s доставлен", userID, orderID)
	return nil
}

func (s *MockNotificationService) SendNewOrderAvailable(orderID string) error {
	log.Printf("[NOTIFY] Новый заказ доступен: %s", orderID)
	return nil
}
