package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr_reviewer/internal/models"
	"pr_reviewer/internal/storage"
	"strconv"
)

func SetIsActiveHandler(db storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "user_id required")
			return
		}

		isActive := r.URL.Query().Get("is_active")
		if isActive == "" {
			WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "new status required")
			return
		}

		isActiveBool, err := strconv.ParseBool(isActive)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid is_active value")
			return
		}

		updatedUser, err := db.SetIsActive(userID, isActiveBool)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				WriteError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
				return
			}
			WriteError(w, http.StatusInternalServerError, "SERVER_ERROR", "internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": updatedUser,
		})

	}
}

func GetUserPRsHandler(db storage.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "user_id required")
			return
		}

		prs, err := db.GetPR(userID)
		if err != nil {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
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
