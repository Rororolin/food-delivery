package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// isNotFound — эвристика для разделения «не найдено» и «бизнес-ошибки»
func isNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
