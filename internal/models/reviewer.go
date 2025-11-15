package models

type Reviewer struct {
	PullRequestID string `json:"pull_request_id"`
	ReviewerID    string `json:"reviewer_id"`
}
