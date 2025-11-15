package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pr_reviewer/internal/storage"
)

func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Team created")
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
