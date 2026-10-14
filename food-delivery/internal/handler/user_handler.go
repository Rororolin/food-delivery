package handler

import (
	"bytes"
	"encoding/json"
	"food-delivery/internal/domain"
	"food-delivery/internal/port"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const adminAPI = "http://localhost:8081"

// syncUserToAdmin — отправляет нового пользователя в админский сервис
func syncUserToAdmin(user domain.User) {
	payload := map[string]string{
		"id":    user.ID,
		"name":  user.Name,
		"phone": user.Phone,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", adminAPI+"/admin/users/sync", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 3 * 1000000000} // 3s
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}

type UserHandler struct {
	repo port.UserRepository
}

func NewUserHandler(repo port.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// POST /users/register — регистрация клиента при оформлении заказа
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	// Если пользователь уже существует — не перезаписываем,
	// чтобы не сбросить IsBlocked и другие поля
	if existing, err := h.repo.GetByID(body.ID); err == nil {
		writeJSON(w, http.StatusOK, existing)
		return
	}

	user := domain.User{
		ID:        body.ID,
		Name:      body.Name,
		Phone:     body.Phone,
		CreatedAt: time.Now(),
	}
	_ = h.repo.Save(user)
	go syncUserToAdmin(user)
	writeJSON(w, http.StatusOK, user)
}

// PATCH /users/{user_id}/block-status — синхронизация статуса блокировки от админа
func (h *UserHandler) UpdateBlockStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	var body struct {
		IsBlocked bool `json:"is_blocked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetByID(userID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}
	user.IsBlocked = body.IsBlocked
	_ = h.repo.Save(user)
	writeJSON(w, http.StatusOK, map[string]any{"user_id": userID, "is_blocked": user.IsBlocked})
}
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, _ := h.repo.GetAll()
	writeJSON(w, http.StatusOK, map[string]any{"users": users})
}
