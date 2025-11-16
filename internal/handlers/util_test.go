package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pr_reviewer/internal/handlers"
	"pr_reviewer/internal/models"
)

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	handlers.WriteError(rec, http.StatusBadRequest, "NOT_FOUND", "user not found")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	var resp models.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected code NOT_FOUND, got %s", resp.Error.Code)
	}
	if resp.Error.Message != "user not found" {
		t.Errorf("expected message 'user not found', got %s", resp.Error.Message)
	}
}
