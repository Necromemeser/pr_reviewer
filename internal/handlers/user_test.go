package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pr_reviewer/internal/handlers"
	"pr_reviewer/internal/models"
	"pr_reviewer/internal/storage"
	mocks "pr_reviewer/internal/storage/mocks"

	"github.com/golang/mock/gomock"
)

func TestSetIsActiveHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockUserStorage(ctrl)

	handler := handlers.SetIsActiveHandler(mockDB)

	userID := "u1"
	newStatus := true

	mockDB.
		EXPECT().
		SetIsActive(userID, newStatus).
		Return(&models.User{
			UserID:   userID,
			Username: "Ivan",
			IsActive: newStatus,
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/user/set?user_id=u1&is_active=true", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]models.User
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["user"].IsActive != true {
		t.Fatalf("expected IsActive true, got %v", resp["user"].IsActive)
	}
}

func TestSetIsActiveHandler_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockUserStorage(ctrl)

	handler := handlers.SetIsActiveHandler(mockDB)

	mockDB.
		EXPECT().
		SetIsActive("u2", true).
		Return(nil, storage.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/user/set?user_id=u2&is_active=true", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestGetUserPRsHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockUserStorage(ctrl)

	handler := handlers.GetUserPRsHandler(mockDB)

	userID := "u1"

	mockDB.
		EXPECT().
		GetPR(userID).
		Return([]models.PullRequestShort{
			{PullRequestID: "pr-1001", PullRequestName: "Add search"},
			{PullRequestID: "pr-1002", PullRequestName: "Fix bug"},
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/user/prs?user_id=u1", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		UserID       string                    `json:"user_id"`
		PullRequests []models.PullRequestShort `json:"pull_requests"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.PullRequests) != 2 {
		t.Fatalf("expected 2 PRs, got %d", len(resp.PullRequests))
	}
}
