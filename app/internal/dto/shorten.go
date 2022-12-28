package dto

import (
	"cc/app/internal/apperror"
	"cc/app/pkg/base62"
	"cc/app/pkg/urlutils"
	"time"
	"unicode/utf8"
)

type CreateShorten struct {
	ID      string `json:"id"`
	LongURL string `json:"long_url"`
	Title   string `json:"title"`
}

type UpdateShorten struct {
	Title string `json:"title"`
}

type GetShortenStats struct {
	From  string `form:"from"`
	To    string `form:"to"`
	Unit  string `form:"unit"`
	Units int    `form:"units"`
}

type GetShortenSummaryStats struct {
	From string `form:"from"`
	To   string `form:"to"`
}

func (createShorten CreateShorten) Validate() error {
	if createShorten.LongURL == "" {
		return apperror.ErrInvalidParams.SetMessage("long_url is required")
	}

	if err := urlutils.Validate(createShorten.LongURL); err != nil {
		return apperror.ErrInvalidParams.SetError(err).SetMessage("long_url is invalid")
	}

	if _, err := base62.Decode(createShorten.ID); err != nil {
		return apperror.ErrInvalidParams.SetError(err).SetMessage("id is invalid")
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
	case "minute", "hour", "day", "week", "month", "year":
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
