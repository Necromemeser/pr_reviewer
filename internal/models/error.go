package models

type ErrorCode string

const (
	TEAM_EXISTS     ErrorCode = "TEAM_EXISTS"
	PR_EXISTS       ErrorCode = "PR_EXISTS"
	PR_MERGED       ErrorCode = "PR_MERGED"
	NOT_ASSIGNED    ErrorCode = "NOT_ASSIGNED"
	NO_CANDIDATE    ErrorCode = "NO_CANDIDATE"
	NOT_FOUND       ErrorCode = "NOT_FOUND"
	USER_NOT_ACTIVE ErrorCode = "USER_NOT_ACTIVE"
	BAD_REQUEST     ErrorCode = "BAD_REQUEST"
	SERVER_ERROR    ErrorCode = "SERVER_ERROR"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}
