package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestCourierOrderWorkflow — полный воркфлоу курьера:
// зарегистрироваться → добавить заказ в пул → получить доступные заказы → принять → завершить
func TestCourierOrderWorkflow(t *testing.T) {
	base := "http://localhost:8082"

	// Шаг 1: зарегистрировать курьера
	courierBody := map[string]any{
		"name":  "Тест Курьеров",
		"phone": "+79001234567",
	}
	data, _ := json.Marshal(courierBody)
	resp, err := http.Post(base+"/courier/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 1 провален: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("Шаг 1: ожидали 201, получили %d", resp.StatusCode)
	}
	var courierResp map[string]any
	json.NewDecoder(resp.Body).Decode(&courierResp)
	courierID := courierResp["id"].(string)
	t.Logf("✓ Шаг 1: курьер зарегистрирован, ID = %s", courierID)

	// Шаг 2: добавить заказ в пул доступных
	orderBody := map[string]any{
		"id":                  "integration-courier-order-001",
		"user_id":             "u1",
		"customer_name":       "Тестовый Клиент",
		"restaurant_name":     "Суши Токио",
		"restaurant_address":  "Невский проспект, 50",
		"restaurant_location": map[string]float64{"lat": 59.934, "lon": 30.335},
		"delivery_address":    "Московский проспект, 15",
		"delivery_location":   map[string]float64{"lat": 59.920, "lon": 30.320},
		"items_count":         2,
		"delivery_fee":        150,
		"total_amount":        800,
	}
	data, _ = json.Marshal(orderBody)
	resp, err = http.Post(base+"/courier/orders", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 2 провален: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("Шаг 2: ожидали 201, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 2: заказ добавлен в пул доступных")

	// Шаг 3: получить список доступных заказов
	resp, err = http.Get(base + "/courier/orders/available")
	if err != nil {
		t.Fatalf("Шаг 3 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 3: ожидали 200, получили %d", resp.StatusCode)
	}
	var availResp map[string]any
	json.NewDecoder(resp.Body).Decode(&availResp)
	orders := availResp["orders"].([]any)
	if len(orders) == 0 {
		t.Fatal("Шаг 3: список доступных заказов пуст")
	}
	t.Logf("✓ Шаг 3: доступных заказов в пуле: %d", len(orders))

	// Шаг 4: курьер принимает заказ
	acceptBody := map[string]any{
		"courier_id": courierID,
		"action":     "accept",
	}
	data, _ = json.Marshal(acceptBody)
	resp, err = http.Post(
		base+"/courier/orders/integration-courier-order-001/respond",
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 4 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 4: ожидали 200, получили %d", resp.StatusCode)
	}
	var acceptResp map[string]any
	json.NewDecoder(resp.Body).Decode(&acceptResp)
	t.Logf("✓ Шаг 4: заказ принят, статус = '%v'", acceptResp["status"])

	// Шаг 5: попытка занятого курьера взять второй заказ (логическая ошибка)
	// Сначала добавим второй заказ в пул
	orderBody2 := map[string]any{
		"id":                  "integration-courier-order-002",
		"user_id":             "u2",
		"customer_name":       "Второй Клиент",
		"restaurant_name":     "Пицца Рома",
		"restaurant_address":  "Садовая улица, 10",
		"restaurant_location": map[string]float64{"lat": 59.928, "lon": 30.320},
		"delivery_address":    "Лиговский проспект, 25",
		"delivery_location":   map[string]float64{"lat": 59.921, "lon": 30.354},
		"items_count":         1,
		"delivery_fee":        150,
		"total_amount":        520,
	}
	data, _ = json.Marshal(orderBody2)
	http.Post(base+"/courier/orders", "application/json", bytes.NewBuffer(data))

	// Пробуем принять второй заказ — должен вернуть 409
	data, _ = json.Marshal(acceptBody)
	resp, err = http.Post(
		base+"/courier/orders/integration-courier-order-002/respond",
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 5 провален: %v", err)
	}
	if resp.StatusCode != 409 {
		t.Fatalf("Шаг 5: ожидали 409, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 5: занятый курьер не может взять второй заказ (409)")

	// Шаг 6: получить маршрут доставки
	resp, err = http.Get(
		base + "/courier/orders/integration-courier-order-001/route?courier_id=" + courierID)
	if err != nil {
		t.Fatalf("Шаг 6 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 6: ожидали 200, получили %d", resp.StatusCode)
	}
	var routeResp map[string]any
	json.NewDecoder(resp.Body).Decode(&routeResp)
	t.Logf("✓ Шаг 6: маршрут построен, расстояние = %.2f км, время = %.0f мин",
		routeResp["total_distance_km"],
		routeResp["total_estimated_minutes"])

	// Шаг 7: подтвердить доставку
	completeBody := map[string]any{"courier_id": courierID}
	data, _ = json.Marshal(completeBody)
	resp, err = http.Post(
		base+"/courier/orders/integration-courier-order-001/complete",
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 7 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 7: ожидали 200, получили %d", resp.StatusCode)
	}
	var completeResp map[string]any
	json.NewDecoder(resp.Body).Decode(&completeResp)
	t.Logf("✓ Шаг 7: доставка подтверждена, статус заказа = '%v'", completeResp["status"])

	// Шаг 8: проверить что курьер снова свободен
	resp, err = http.Get(base + "/courier/profile?courier_id=" + courierID)
	if err != nil {
		t.Fatalf("Шаг 8 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 8: ожидали 200, получили %d", resp.StatusCode)
	}
	var profileResp map[string]any
	json.NewDecoder(resp.Body).Decode(&profileResp)
	if profileResp["status"] != "free" {
		t.Fatalf("Шаг 8: ожидали статус курьера 'free', получили '%v'", profileResp["status"])
	}
	t.Logf("✓ Шаг 8: курьер освобождён, статус = '%v'", profileResp["status"])
}
