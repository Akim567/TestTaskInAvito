package domain

type ErrorCode string

const (
	ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrorCodePRExists    ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged    ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"
)

type DomainError struct {
	Code    ErrorCode
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewTeamExistsError(teamName string) *DomainError {
	return &DomainError{
		Code:    ErrorCodeTeamExists,
		Message: "team already exists: " + teamName,
	}
}

func NewPRExistsError(prID string) *DomainError {
	return &DomainError{
		Code:    ErrorCodePRExists,
		Message: "pull request already exists: " + prID,
	}
}

func NewPRMergedError(prID string) *DomainError {
	return &DomainError{
		Code:    ErrorCodePRMerged,
		Message: "pull request already merged: " + prID,
	}
}

func NewNotAssignedError(userID, prID string) *DomainError {
	return &DomainError{
		Code:    ErrorCodeNotAssigned,
		Message: "user " + userID + " is not assigned to PR " + prID,
	}
}

func NewNoCandidateError(teamName string) *DomainError {
	return &DomainError{
		Code:    ErrorCodeNoCandidate,
		Message: "no active candidate in team: " + teamName,
	}
}

func NewNotFoundError(resource string) *DomainError {
	return &DomainError{
		Code:    ErrorCodeNotFound,
		Message: resource + " not found",
	}
}
