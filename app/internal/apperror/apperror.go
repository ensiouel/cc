package apperror

import (
	"cc/app/pkg/errs"
)

const (
	TypeUnknown errs.Type = iota
	TypeInternal
	TypeNotExists
	TypeAlreadyExists
	TypeInvalidParams
	TypeInvalidCredentials
	TypeUnauthorized
	TypeNotOwned
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
	Unknown            = TypeUnknown.New("unknown error")
	Internal           = TypeInternal.New("internal error")
	NotExists          = TypeNotExists.New("not exists")
	AlreadyExists      = TypeAlreadyExists.New("already exists")
	InvalidParams      = TypeInvalidParams.New("invalid params")
	InvalidCredentials = TypeInvalidCredentials.New("invalid credentials")
	Unauthorized       = TypeUnauthorized.New("unauthorized")
	NotOwned           = TypeNotOwned.New("not owned")
)
