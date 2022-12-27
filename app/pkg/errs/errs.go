package errs

import (
	"fmt"
)

type ErrorCode int

type Error struct {
	Code    ErrorCode `json:"code"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Scope   string    `json:"-"`
	Err     error     `json:"-"`
}

func (code ErrorCode) New(status string) *Error {
	return &Error{
		Code:    code,
		Status:  status,
		Message: "",
		Err:     nil,
	}
}

func (error Error) Is(target error) bool {
	err, ok := target.(*Error)
	if !ok {
		return false
	}

	return error.Code == err.Code
}

func (error Error) Error() string {
	if error.Err != nil {
		if error.Scope == "" {
			return fmt.Sprintf("%s: %s", error.Status, error.Err.Error())
		}

		return fmt.Sprintf("%s: %s: %s", error.Status, error.Scope, error.Err.Error())
	}

	if error.Scope == "" {
		return fmt.Sprintf("%s: %s", error.Status, error.Message)
	}

	return fmt.Sprintf("%s: %s: %s", error.Status, error.Scope, error.Message)
}

func (error Error) SetError(err error) Error {
	if error.Err != nil {
		error.Err = fmt.Errorf("%s: %w", error.Err, err)
	} else {
		error.Err = err
	}

	return error
}

func (error Error) SetScope(scope string) Error {
	error.Scope = scope

	return error
}

func (error Error) SetMessage(message string) Error {
	error.Message = message

	return error
}
