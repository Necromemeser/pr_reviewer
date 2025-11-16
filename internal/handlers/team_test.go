package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestDeactivateTeamHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTeamStorage := mocks.NewMockTeamStorage(ctrl)
	handler := handlers.DeactivateTeamHandler(mockTeamStorage)

	doReq := func(url string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, url, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	// team_name missing 400 BAD_REQUEST
	rec := doReq("/team/deactivate")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	// team not found 404 NOT_FOUND
	mockTeamStorage.EXPECT().GetTeam("backend").Return(nil, storage.ErrTeamNotFound)

	rec = doReq("/team/deactivate?team_name=backend")
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}

	// unexpected GetTeam error 500 SERVER_ERROR
	mockTeamStorage.EXPECT().GetTeam("backend").Return(nil, errors.New("db error"))

	rec = doReq("/team/deactivate?team_name=backend")
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}

	// DeactivateTeam error 500 SERVER_ERROR
	team := &models.Team{TeamName: "backend"}

	mockTeamStorage.EXPECT().GetTeam("backend").Return(team, nil)
	mockTeamStorage.EXPECT().DeactivateTeam(*team).Return(nil, errors.New("update error"))

	rec = doReq("/team/deactivate?team_name=backend")
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}

	// 5. SUCCESS 200 OK
	deactivated := &models.Team{TeamName: "backend"}

	mockTeamStorage.EXPECT().GetTeam("backend").Return(team, nil)
	mockTeamStorage.EXPECT().DeactivateTeam(*team).Return(deactivated, nil)

	rec = doReq("/team/deactivate?team_name=backend")
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
