package handlers_test

import (
	"bytes"
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

func TestCreateTeamHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockTeamStorage(ctrl)

	handler := handlers.CreateTeamHandler(mockDB)

	team := models.Team{
		TeamName: "backend",
	}

	mockDB.
		EXPECT().
		AddTeam(team).
		Return(&team, nil)

	body, _ := json.Marshal(team)

	req := httptest.NewRequest(http.MethodPost, "/team", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp map[string]models.Team
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["team"].TeamName != "backend" {
		t.Fatalf("unexpected response: %v", resp)
	}
}

func TestCreateTeamHandler_TeamExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockTeamStorage(ctrl)
	handler := handlers.CreateTeamHandler(mockDB)

	team := models.Team{TeamName: "backend"}

	mockDB.
		EXPECT().
		AddTeam(team).
		Return(nil, storage.ErrTeamExists)

	body, _ := json.Marshal(team)

	req := httptest.NewRequest("POST", "/team", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetTeamHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockTeamStorage(ctrl)
	handler := handlers.GetTeamHandler(mockDB)

	team := models.Team{TeamName: "frontend"}

	mockDB.
		EXPECT().
		GetTeam("frontend").
		Return(&team, nil)

	req := httptest.NewRequest("GET", "/team?team_name=frontend", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result models.Team
	json.NewDecoder(w.Body).Decode(&result)

	if result.TeamName != "frontend" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestGetTeamHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockTeamStorage(ctrl)
	handler := handlers.GetTeamHandler(mockDB)

	mockDB.
		EXPECT().
		GetTeam("unknown").
		Return(nil, storage.ErrTeamNotFound)

	req := httptest.NewRequest("GET", "/team?team_name=unknown", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
