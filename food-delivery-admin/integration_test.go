package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestAdminOrderWorkflow — полный воркфлоу администратора:
// создать заказ → изменить статус → назначить курьера → отменить → возврат средств
func TestAdminOrderWorkflow(t *testing.T) {
	base := "http://localhost:8081"

	// Шаг 1: создать заказ через административный сервис
	orderBody := map[string]any{
		"id":          "integration-test-order-001",
		"customer_id": "Тестовый Клиент · +79001234567",
		"totalAmount": 1500.00,
		"is_paid":     true,
	}
	data, _ := json.Marshal(orderBody)
	resp, err := http.Post(base+"/admin/orders", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Шаг 1 провален: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Fatalf("Шаг 1: ожидали 201, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 1: заказ создан в административной панели")

	// Шаг 2: получить список заказов и убедиться что заказ появился
	resp, err = http.Get(base + "/admin/orders")
	if err != nil {
		t.Fatalf("Шаг 2 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 2: ожидали 200, получили %d", resp.StatusCode)
	}
	var listResp map[string]any
	json.NewDecoder(resp.Body).Decode(&listResp)
	orders := listResp["orders"].([]any)
	if len(orders) == 0 {
		t.Fatal("Шаг 2: список заказов пуст — заказ не сохранился")
	}
	t.Logf("✓ Шаг 2: список заказов получен, заказов в системе: %d", len(orders))

	// Шаг 3: изменить статус заказа с pending на assigned
	statusBody := map[string]any{"status": "assigned"}
	data, _ = json.Marshal(statusBody)
	req, _ := http.NewRequest("PATCH",
		base+"/admin/orders/integration-test-order-001/status",
		bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Шаг 3 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 3: ожидали 200, получили %d", resp.StatusCode)
	}
	var statusResp map[string]any
	json.NewDecoder(resp.Body).Decode(&statusResp)
	t.Logf("✓ Шаг 3: статус изменён на '%v'", statusResp["new_status"])

	// Шаг 4: попытка недопустимого перехода статуса (assigned → pending)
	invalidBody := map[string]any{"status": "pending"}
	data, _ = json.Marshal(invalidBody)
	req, _ = http.NewRequest("PATCH",
		base+"/admin/orders/integration-test-order-001/status",
		bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Шаг 4 провален: %v", err)
	}
	if resp.StatusCode != 409 {
		t.Fatalf("Шаг 4: ожидали 409, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 4: недопустимый переход статуса корректно отклонён (409)")

	// Шаг 5: отменить заказ
	req, _ = http.NewRequest("POST",
		base+"/admin/orders/integration-test-order-001/cancel",
		nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Шаг 5 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 5: ожидали 200, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 5: заказ отменён")

	// Шаг 6: выполнить возврат средств по отменённому оплаченному заказу
	req, _ = http.NewRequest("POST",
		base+"/admin/orders/integration-test-order-001/refund",
		nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Шаг 6 провален: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Шаг 6: ожидали 200, получили %d", resp.StatusCode)
	}
	var refundResp map[string]any
	json.NewDecoder(resp.Body).Decode(&refundResp)
	t.Logf("✓ Шаг 6: возврат выполнен, статус возврата = '%v'", refundResp["refund_status"])

	// Шаг 7: повторный возврат — должен вернуть 409
	req, _ = http.NewRequest("POST",
		base+"/admin/orders/integration-test-order-001/refund",
		nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Шаг 7 провален: %v", err)
	}
	if resp.StatusCode != 409 {
		t.Fatalf("Шаг 7: ожидали 409, получили %d", resp.StatusCode)
	}
	t.Log("✓ Шаг 7: повторный возврат корректно отклонён (409)")
}
