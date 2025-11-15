package handlers

import (
	"net/http"
	"pr_reviewer/internal/storage"
)

func NewRouter(db *storage.Storage) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Team
	// mux.HandleFunc("POST /team/add", CreateTeamHandler(db))
	mux.HandleFunc("GET /team/get", GetTeamHandler(db))

	// PullRequest
	// mux.HandleFunc("POST /pullRequest/create", CreatePullRequestHandler(db))
	// mux.HandleFunc("POST /pullRequest/merge", MergePRHandler(db))
	// mux.HandleFunc("POST /pullRequest/reassign", ReassignReviewerHandler(db))

	// Users
	mux.HandleFunc("GET /users/getReview", GetUserPRsHandler(db))
	mux.HandleFunc("POST /users/setIsActive", SetIsActiveHandler(db))

	return mux
}
