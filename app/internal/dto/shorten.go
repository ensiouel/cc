package dto

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/pkg/base62"
	"cc/app/pkg/urlutils"
	"time"
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

type GetShortenStats struct {
	From  string      `form:"from"`
	To    string      `form:"to"`
	Unit  domain.Unit `form:"unit"`
	Units int         `form:"units"`
}

type ExportShortenStats struct {
	From string `form:"from"`
	To   string `form:"to"`
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

func (getShortenStats GetShortenStats) Validate() error {
	var err error
	_, err = time.Parse("2006-01-02", getShortenStats.From)
	if err != nil {
		return apperror.InvalidParams.WithMessage("from is invalid, expected 2006-01-02")
	}

	_, err = time.Parse("2006-01-02", getShortenStats.To)
	if err != nil {
		return apperror.InvalidParams.WithMessage("to is invalid, expected 2006-01-02")
	}

	switch getShortenStats.Unit {
	case domain.UnitHour, domain.UnitDay, domain.UnitWeek, domain.UnitMonth, domain.UnitYear:
	default:
		return apperror.InvalidParams.WithMessage("unit is invalid, expected (hour, day, week, month, year)")
	}

	return nil
}

func (exportShortenStats ExportShortenStats) Validate() error {
	var err error
	_, err = time.Parse("2006-01-02", exportShortenStats.From)
	if err != nil {
		return apperror.InvalidParams.WithMessage("from is invalid, expected 2006-01-02")
	}

	_, err = time.Parse("2006-01-02", exportShortenStats.To)
	if err != nil {
		return apperror.InvalidParams.WithMessage("to is invalid, expected 2006-01-02")
	}

	return nil
}
