package dto

import (
	"cc/app/internal/apperror"
	"github.com/google/uuid"
	"unicode/utf8"
)

type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (credentials Credentials) Validate() error {
	if credentials.Name == "" {
		return apperror.InvalidParams.WithMessage("name is required")
	}

	if credentials.Password == "" {
		return apperror.InvalidParams.WithMessage("password is required")
	}

	if !utf8.ValidString(credentials.Name) {
		return apperror.InvalidParams.WithMessage("name is invalid")
	}

	if !utf8.ValidString(credentials.Password) {
		return apperror.InvalidParams.WithMessage("password is invalid")
	}

	nameLen := utf8.RuneCountInString(credentials.Name)
	if nameLen < 3 {
		return apperror.InvalidParams.WithMessage("name is to short")
	} else if nameLen > 20 {
		return apperror.InvalidParams.WithMessage("name is to long")
	}

	passwordLen := utf8.RuneCountInString(credentials.Password)
	if passwordLen < 5 {
		return apperror.InvalidParams.WithMessage("password is to short")
	} else if passwordLen > 50 {
		return apperror.InvalidParams.WithMessage("password is to long")
	}

	return nil
}

type Refresh struct {
	RefreshToken uuid.UUID `json:"refresh_token"`
}
