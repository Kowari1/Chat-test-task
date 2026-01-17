package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func writeJSON(w http.ResponseWriter, status int, data any, logger *zap.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to encode JSON response", zap.Error(err))
		http.Error(w, ErrInternal, http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, msg string, logger *zap.Logger) {
	writeJSON(w, status, map[string]string{"error": msg}, logger)
}
