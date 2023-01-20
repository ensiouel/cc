package apperror

import (
	"cc/app/pkg/errs"
)

const (
	Unknown errs.ErrorCode = iota - 1
	Internal
	NotExists
	AlreadyExists
	InvalidParams
	InvalidCredentials
	Unauthorized
)

func Is(target error, code errs.ErrorCode) (err errs.Error, ok bool) {
	err, ok = target.(errs.Error)
	if !ok {
		return
	}

	if err.Code != code {
		return err, false
	}

	return
}

var (
	ErrUnknownError       = Unknown.New("unknown error")
	ErrInternalError      = Internal.New("internal error")
	ErrNotExists          = NotExists.New("not exists")
	ErrAlreadyExists      = AlreadyExists.New("already exists")
	ErrInvalidParams      = InvalidParams.New("invalid params")
	ErrInvalidCredentials = InvalidCredentials.New("invalid credentials")
	ErrUnauthorized       = Unauthorized.New("unauthorized")
)
