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
			http.Error(w, "team_name required", http.StatusBadRequest)
			return
		}

		team, err := db.GetTeam(teamName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(team)
	}
}
