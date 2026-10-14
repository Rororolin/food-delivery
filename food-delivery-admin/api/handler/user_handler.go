package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"food-delivery-admin/internal/usecase"

	"github.com/go-chi/chi/v5"
)

const clientAPI = "http://localhost:8080"

// syncBlockToClient — синхронизирует статус блокировки с клиентским сервисом
func syncBlockToClient(userID string, blocked bool) {
	payload := map[string]interface{}{"is_blocked": blocked}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PATCH", clientAPI+"/users/"+userID+"/block-status", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}

// UserHandler — REST-обработчики для Сервиса управления пользователями
type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// SyncUser POST /admin/users/sync — приём пользователя от клиентского сервиса
func (h *UserHandler) SyncUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	// Если пользователь уже есть — не перезаписываем (чтобы не сбросить IsBlocked)
	if existing, err := h.uc.GetUserByID(body.ID); err == nil && existing != nil {
		writeJSON(w, http.StatusOK, existing)
		return
	}

	user, err := h.uc.CreateUser(body.ID, body.Name, body.Phone)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}


func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"users": users})
}

// GetUser GET /admin/users/{user_id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	user, err := h.uc.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// BlockUser PATCH /admin/users/{user_id}/block
func (h *UserHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	user, err := h.uc.BlockUser(userID)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	go syncBlockToClient(userID, true)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": user.ID,
		"blocked": user.IsBlocked,
	})
}

// UnblockUser PATCH /admin/users/{user_id}/unblock
func (h *UserHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	user, err := h.uc.UnblockUser(userID)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusConflict, err.Error())
		}
		return
	}

	go syncBlockToClient(userID, false)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": user.ID,
		"blocked": user.IsBlocked,
	})
}
