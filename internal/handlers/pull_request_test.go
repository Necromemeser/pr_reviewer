package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pr_reviewer/internal/handlers"
	"pr_reviewer/internal/models"
	mocks "pr_reviewer/internal/storage/mocks"

	"github.com/golang/mock/gomock"
)

func TestCreatePullRequestHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPRStorage(ctrl)
	handler := handlers.CreatePullRequestHandler(mockDB)

	authorID := "u1"
	pr := models.PullRequest{
		PullRequestID:   "pr-1001",
		PullRequestName: "Add search",
		AuthorID:        &authorID,
	}

	mockDB.
		EXPECT().
		CreatePullRequest(pr.PullRequestID, pr.PullRequestName, *pr.AuthorID).
		Return(&pr, nil)

	body, _ := json.Marshal(pr)
	req := httptest.NewRequest(http.MethodPost, "/pr", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp map[string]models.PullRequest
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["pr"].PullRequestID != pr.PullRequestID {
		t.Fatalf("unexpected PR in response: %v", resp["pr"])
	}
}

func TestMergePRHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPRStorage(ctrl)
	handler := handlers.MergePRHandler(mockDB)

	prID := "pr-1001"
	mergedPR := models.PullRequest{
		PullRequestID:   prID,
		PullRequestName: "Add search",
		Status:          "MERGED",
	}

	mockDB.
		EXPECT().
		MergePullRequest(prID).
		Return(&mergedPR, nil)

	body, _ := json.Marshal(map[string]string{"pull_request_id": prID})
	req := httptest.NewRequest(http.MethodPost, "/pr/merge", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp map[string]models.PullRequest
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["pr"].Status != "MERGED" {
		t.Fatalf("expected status MERGED, got %s", resp["pr"].Status)
	}
}

func TestReassignReviewerHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPRStorage(ctrl)
	handler := handlers.ReassignReviewerHandler(mockDB)

	prID := "pr-1001"
	oldReviewer := "u2"
	newReviewer := "u3"

	newPR := &models.PullRequest{
		PullRequestID:   prID,
		PullRequestName: "Add search",
		Status:          "OPEN",
	}

	mockDB.
		EXPECT().
		ReasignReviewer(prID, oldReviewer).
		Return(newPR, newReviewer, nil)

	body, _ := json.Marshal(map[string]string{
		"pull_request_id": prID,
		"reviewer_id":     oldReviewer,
	})
	req := httptest.NewRequest(http.MethodPost, "/pr/reassign", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp struct {
		PR         *models.PullRequest `json:"pr"`
		ReplacedBy string              `json:"replaced_by"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ReplacedBy != newReviewer {
		t.Fatalf("expected replaced_by %s, got %s", newReviewer, resp.ReplacedBy)
	}
	if resp.PR.PullRequestID != prID {
		t.Fatalf("expected PR id %s, got %s", prID, resp.PR.PullRequestID)
	}
}
