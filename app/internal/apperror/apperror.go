package apperror

import (
	"cc/app/pkg/errs"
)

const (
	TypeInternal errs.Type = iota
	TypeNotExists
	TypeAlreadyExists
	TypeInvalidParams
	TypeInvalidCredentials
	TypeUnauthorized
)

func Is(target error, code errs.Type) (err errs.Error, ok bool) {
	err, ok = target.(errs.Error)
	if !ok {
		return
	}

	if err.Type != code {
		return err, false
	}

	return
}

var (
	Internal           = TypeInternal.New("internal error")
	NotExists          = TypeNotExists.New("not exists")
	AlreadyExists      = TypeAlreadyExists.New("already exists")
	InvalidParams      = TypeInvalidParams.New("invalid params")
	InvalidCredentials = TypeInvalidCredentials.New("invalid credentials")
	Unauthorized       = TypeUnauthorized.New("unauthorized")
)
