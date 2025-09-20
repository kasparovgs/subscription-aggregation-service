package domain

import "fmt"

var (
	CodeNotFound     = 404
	CodeUnauthorized = 401
	CodeAlreadyExist = 409
	CodeForbidden    = 403
)

type MyErr struct {
	Code    int
	Message string
}

func (e *MyErr) Error() string {
	return fmt.Sprintf("[Error %d] %s", e.Code, e.Message)
}

func NewError(code int, msg string) *MyErr {
	return &MyErr{
		Code:    code,
		Message: msg,
	}
}

var (
	ErrNotFound = func(msg string) *MyErr {
		return NewError(CodeNotFound, "Not found: "+msg)
	}
	ErrUnauthorized = func(msg string) *MyErr {
		return NewError(CodeUnauthorized, "Unauthorized: "+msg)
	}
	ErrAlreadyExist = func(msg string) *MyErr {
		return NewError(CodeAlreadyExist, "Already exist: "+msg)
	}
	ErrForbidden = func(msg string) *MyErr {
		return NewError(CodeForbidden, "Forbidden: "+msg)
	}
)
