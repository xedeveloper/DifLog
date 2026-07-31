package errors

import "fmt"

type ErrorCode string

const (
	ErrNotInitialized    ErrorCode = "NOT_INITIALIZED"
	ErrAlreadyInitialized ErrorCode = "ALREADY_INITIALIZED"
	ErrNoStagedContexts  ErrorCode = "NO_STAGED_CONTEXTS"
	ErrContextNotFound   ErrorCode = "CONTEXT_NOT_FOUND"
	ErrCommitNotFound    ErrorCode = "COMMIT_NOT_FOUND"
	ErrBranchNotFound    ErrorCode = "BRANCH_NOT_FOUND"
	ErrBranchExists      ErrorCode = "BRANCH_EXISTS"
	ErrRemoteNotSet      ErrorCode = "REMOTE_NOT_SET"
	ErrAIToolNotDetected ErrorCode = "AI_TOOL_NOT_DETECTED"
	ErrStorageFailure    ErrorCode = "STORAGE_FAILURE"
	ErrHashingFailure    ErrorCode = "HASHING_FAILURE"
	ErrRemoteFailure     ErrorCode = "REMOTE_FAILURE"
	ErrInvalidArgument   ErrorCode = "INVALID_ARGUMENT"
	ErrUnsupportedPlatform ErrorCode = "UNSUPPORTED_PLATFORM"
)

type DifLogError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *DifLogError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DifLogError) Unwrap() error {
	return e.Err
}

func NewDifLogError(code ErrorCode, message string, err error) *DifLogError {
	return &DifLogError{Code: code, Message: message, Err: err}
}
