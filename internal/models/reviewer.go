package models

type Reviewer struct {
	PullRequestID string `json:"pull_request_id"`
	ReviewerID    string `json:"old_user_id"`
}
