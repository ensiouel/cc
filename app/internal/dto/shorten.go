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
	Title string `json:"title"`
	URL   string `json:"url"`
}

type GetShortenStats struct {
	From  string      `form:"from"`
	To    string      `form:"to"`
	Unit  domain.Unit `form:"unit"`
	Units int         `form:"units"`
}

type GetShortenSummaryStats struct {
	From string `form:"from"`
	To   string `form:"to"`
}

func (createShorten CreateShorten) Validate() error {
	if createShorten.URL == "" {
		return apperror.ErrInvalidParams.SetMessage("url is required")
	}

	if err := urlutils.Validate(createShorten.URL); err != nil {
		return apperror.ErrInvalidParams.SetError(err).SetMessage("url is invalid")
	}

	if _, err := base62.Decode(createShorten.Key); err != nil {
		return apperror.ErrInvalidParams.SetError(err).SetMessage("key is invalid")
	}

	return nil
}

func (updateShorten UpdateShorten) Validate() error {
	if updateShorten.Title == "" {
		return apperror.ErrInvalidParams.SetMessage("title is required")
	}

	if utf8.RuneCountInString(updateShorten.Title) > 100 {
		return apperror.ErrInvalidParams.SetMessage("title is to long")
	}

	if err := urlutils.Validate(updateShorten.URL); err != nil {
		return apperror.ErrInvalidParams.SetError(err).SetMessage("url is invalid")
	}

	return nil
}

func (getShortenStats GetShortenStats) Validate() error {
	var err error
	_, err = time.Parse("2006-01-02", getShortenStats.From)
	if err != nil {
		return apperror.ErrInvalidParams.SetMessage("from is invalid, expected 2006-01-02")
	}

	_, err = time.Parse("2006-01-02", getShortenStats.To)
	if err != nil {
		return apperror.ErrInvalidParams.SetMessage("to is invalid, expected 2006-01-02")
	}

	switch getShortenStats.Unit {
	case domain.UnitMinute, domain.UnitHour, domain.UnitDay, domain.UnitWeek, domain.UnitMonth, domain.UnitYear:
	default:
		return apperror.ErrInvalidParams.SetMessage("unit is invalid, expected (minute, hour, day, week, month, year)")
	}

	return nil
}

func (getShortenSummaryStats GetShortenSummaryStats) Validate() error {
	var err error
	_, err = time.Parse("2006-01-02", getShortenSummaryStats.From)
	if err != nil {
		return apperror.ErrInvalidParams.SetMessage("from is invalid, expected 2006-01-02")
	}

	_, err = time.Parse("2006-01-02", getShortenSummaryStats.To)
	if err != nil {
		return apperror.ErrInvalidParams.SetMessage("to is invalid, expected 2006-01-02")
	}

	return nil
}
