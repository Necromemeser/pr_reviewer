package handlers

import (
	"encoding/json"
	"net/http"
	"pr_reviewer/internal/models"
)

func WriteError(w http.ResponseWriter, status int, code models.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
