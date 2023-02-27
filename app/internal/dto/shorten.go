package dto

import (
	"cc/app/internal/apperror"
	"cc/app/pkg/base62"
	"cc/app/pkg/urlutils"
	"unicode/utf8"
)

type CreateShorten struct {
	Key   string `json:"key"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type UpdateShorten struct {
	Title string   `json:"title,omitempty"`
	URL   string   `json:"url,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

type SelectShortens struct {
	Tags []string `form:"tags,omitempty"`
}

func (createShorten CreateShorten) Validate() error {
	if createShorten.URL == "" {
		return apperror.InvalidParams.WithMessage("url is required")
	}

	if err := urlutils.Validate(createShorten.URL); err != nil {
		return apperror.InvalidParams.WithError(err).WithMessage("url is invalid")
	}

	if _, err := base62.Decode(createShorten.Key); err != nil {
		return apperror.InvalidParams.WithError(err).WithMessage("key is invalid")
	}

	return nil
}

func (updateShorten UpdateShorten) Validate() error {
	if updateShorten.Title != "" && utf8.RuneCountInString(updateShorten.Title) > 100 {
		return apperror.InvalidParams.WithMessage("title is to long")
	}

	if updateShorten.URL != "" && urlutils.Validate(updateShorten.URL) != nil {
		return apperror.InvalidParams.WithMessage("url is invalid")

	}

	return nil
}
