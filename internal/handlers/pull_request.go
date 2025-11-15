package handlers

import (
	"encoding/json"
	"net/http"
	"pr_reviewer/internal/models"
	"pr_reviewer/internal/storage"
)

func CreatePullRequestHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var pr models.PullRequest
		if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
			return
		}

		if pr.AuthorID == nil || *pr.AuthorID == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "author_id is required")
			return
		}
		if pr.PullRequestID == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "pull_request_id is required")
			return
		}
		if pr.PullRequestName == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "pull_request_name is required")
			return
		}

		newPR, err := db.CreatePullRequest(pr.PullRequestID, pr.PullRequestName, *pr.AuthorID)
		if err == storage.ErrPRExists {
			WriteError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
			return
		}
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "NOT_FOUND", "unexpected error")
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"pr": newPR,
		})

	}
}

func MergePRHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var pr models.PullRequest
		if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
			return
		}

		if pr.PullRequestID == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "pull_request_id is required")
			return
		}

		mergedPR, err := db.MergePullRequest(pr.PullRequestID)
		if err == storage.ErrPRNotFound {
			WriteError(w, http.StatusNotFound, "NOT FOUND", "PR id not found")
			return
		}
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "NOT_FOUND", "unexpected error")
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"pr": mergedPR,
		})

	}
}

func ReassignReviewerHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var reviewer models.Reviewer
		if err := json.NewDecoder(r.Body).Decode(&reviewer); err != nil {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
			return
		}

		if reviewer.PullRequestID == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "pull_request_id is required")
			return
		}

		if reviewer.ReviewerID == "" {
			WriteError(w, http.StatusBadRequest, "NOT_FOUND", "author_id is required")
			return
		}

		newPR, newReviewer, err := db.ReasignReviewer(reviewer.PullRequestID, reviewer.ReviewerID)
		if err == storage.ErrPRNotFound {
			WriteError(w, http.StatusNotFound, "NOT FOUND", "PR id not found")
			return
		}
		if err == storage.ErrReviewerNotFound {
			WriteError(w, http.StatusNotFound, "NOT FOUND", "reviewer id not found")
			return
		}
		if err == storage.ErrNoCandidate {
			WriteError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
			return
		}
		if err == storage.ErrPRMerged {
			WriteError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
			return
		}
		if err == storage.ErrPRMerged {
			WriteError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
			return
		}
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "NOT_FOUND", "internal server error")
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"pr":          newPR,
			"replaced_by": newReviewer,
		})

	}
}
