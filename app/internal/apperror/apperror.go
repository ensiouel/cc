package apperror

import (
	"cc/app/pkg/errs"
)

const (
	codeUnknownError errs.ErrorCode = iota - 1
	codeInternalError
	codeNotFound
	codeAlreadyExists
	codeInvalidParams
	codeInvalidCredentials
	codeUnauthorized
)

func Internal(err error) (internal errs.Error, ok bool) {
	internal, ok = err.(errs.Error)
	if !ok {
		return
	}

	if internal.Code != codeInternalError {
		return internal, false
	}

	return
}

var (
	ErrUnknownError       = codeUnknownError.New("unknown error")
	ErrInternalError      = codeInternalError.New("internal error")
	ErrNotExists          = codeNotFound.New("not exists")
	ErrAlreadyExists      = codeAlreadyExists.New("already exists")
	ErrInvalidParams      = codeInvalidParams.New("invalid params")
	ErrInvalidCredentials = codeInvalidCredentials.New("invalid credentials")
	ErrUnauthorized       = codeUnauthorized.New("unauthorized")
)
