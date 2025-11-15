package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr_reviewer/internal/models"
	"pr_reviewer/internal/storage"
	"strconv"
)

func SetIsActiveHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}

		isActive := r.URL.Query().Get("is_active")
		if isActive == "" {
			http.Error(w, "new status required", http.StatusBadRequest)
			return
		}

		isActiveBool, err := strconv.ParseBool(isActive)
		if err != nil {
			http.Error(w, "invalid is_active value", http.StatusBadRequest)
			return
		}

		updatedUser, err := db.SetIsActive(userID, isActiveBool)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": updatedUser,
		})

	}
}

func GetUserPRsHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}

		prs, err := db.GetPR(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		type UserPRsResponse struct {
			UserID       string                    `json:"user_id"`
			PullRequests []models.PullRequestShort `json:"pull_requests"`
		}

		resp := UserPRsResponse{
			UserID:       userID,
			PullRequests: prs,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	}
}
