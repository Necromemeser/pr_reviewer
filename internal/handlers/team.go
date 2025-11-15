package handlers

import (
	"encoding/json"
	"net/http"
	"pr_reviewer/internal/models"
	"pr_reviewer/internal/storage"
)

func CreateTeamHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var team models.Team
		if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
			return
		}

		newTeam, err := db.AddTeam(team)
		if err == storage.ErrTeamExists {
			WriteError(w, http.StatusBadRequest, "TEAM_EXISTS", "team already exists")
			return
		}
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "NOT_FOUND", "unexpected error")
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"team": newTeam,
		})

	}
}

func GetTeamHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "team_name required")
			return
		}

		team, err := db.GetTeam(teamName)
		if err == storage.ErrTeamNotFound {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "team not found")
			return
		}

		if err != nil {
			WriteError(w, http.StatusInternalServerError, "NOT_FOUND", "unexpected error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(team)
	}
}
