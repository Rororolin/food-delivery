package handler

import (
	"encoding/json"
	"food-delivery/internal/port"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CatalogHandler — обработчик компонента управления каталогом
type CatalogHandler struct {
	uc port.CatalogUseCase
}

func NewCatalogHandler(uc port.CatalogUseCase) *CatalogHandler {
	return &CatalogHandler{uc: uc}
}

// GET /restaurants?address=...
func (h *CatalogHandler) ListRestaurants(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, `{"error":"address query param is required"}`, http.StatusBadRequest)
		return
	}

	restaurants, err := h.uc.ListRestaurants(address)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"restaurants": restaurants})
}

// GET /restaurants/{restaurant_id}/menu
func (h *CatalogHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "restaurant_id")

	menu, err := h.uc.GetMenu(restaurantID)
	if err != nil {
		if err.Error() == "restaurant not found" {
			http.Error(w, `{"error":"restaurant not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"menu": menu})
}

// writeJSON — вспомогательная функция для записи JSON-ответа
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
