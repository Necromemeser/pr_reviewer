package models

type ErrorCode string

const (
	TEAM_EXISTS  ErrorCode = "TEAM_EXISTS"
	PR_EXISTS    ErrorCode = "PR_EXISTS"
	PR_MERGED    ErrorCode = "PR_MERGED"
	NOT_ASSIGNED ErrorCode = "NOT_ASSIGNED"
	NO_CANDIDATE ErrorCode = "NO_CANDIDATE"
	NOT_FOUND    ErrorCode = "NOT_FOUND"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}
