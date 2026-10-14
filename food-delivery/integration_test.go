package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestFullOrderWorkflow — полный воркфлоу: рестораны → корзина → заказ → статус
func TestFullOrderWorkflow(t *testing.T) {
	base := "http://localhost:8080"

	// Шаг 1: получить рестораны
	resp, err := http.Get(base + "/restaurants?address=невский")
	if err != nil {
		t.Fatalf("Шаг 1 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 1: ожидали 200, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 1: рестораны получены")

	// Шаг 2: получить меню
	resp, err = http.Get(base + "/restaurants/r1/menu")
	if err != nil {
		t.Fatalf("Шаг 2 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 2: ожидали 200, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 2: меню получено")

	// Шаг 3: добавить в корзину
	cartBody := map[string]any{
		"userId":   "u1",
		"action":   "add",
		"dishId":   "r1m1",
		"quantity": 2,
	}
	data, _ := json.Marshal(cartBody)
	resp, err = http.Post(base+"/cart", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 3 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 3: ожидали 200, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 3: блюдо добавлено в корзину")

	// Шаг 4: оформить заказ
	orderBody := map[string]any{
		"userId": "u1",
		"cart": []map[string]any{
			{"dishId": "r1m1", "quantity": 2},
		},
		"address":       "Невский проспект, 10",
		"paymentMethod": "card",
	}
	data, _ = json.Marshal(orderBody)
	resp, err = http.Post(base+"/orders", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 4 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 4: ожидали 200, получили %d", resp.StatusCode)
	}

	var orderResp map[string]any
	json.NewDecoder(resp.Body).Decode(&orderResp)
	orderID := orderResp["orderId"].(string)
	t.Logf("✓ Шаг 4: заказ создан, ID = %s", orderID)

	// Шаг 5: проверить статус заказа
	resp, err = http.Get(base + "/orders/" + orderID)
	if err != nil {
		t.Fatalf("Шаг 5 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 5: ожидали 200, получили %d", resp.StatusCode)
	}

	var statusResp map[string]any
	json.NewDecoder(resp.Body).Decode(&statusResp)
	if statusResp["status"] != "created" {
		t.Fatalf("Шаг 5: ожидали статус 'created', получили '%v'", statusResp["status"])
	}
	t.Logf("✓ Шаг 5: статус заказа = %v", statusResp["status"])
}
