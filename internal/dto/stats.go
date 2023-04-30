package dto

import (
	"cc/internal/domain"
	"cc/pkg/apperror"
	"time"
)

type CreateClick struct {
	ShortenID uint64    `json:"shorten_id"`
	Platform  string    `json:"platform"`
	OS        string    `json:"os"`
	Referer   string    `json:"referer"`
	IP        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
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

func (getShortenStats GetShortenStats) Validate() error {
	var err error
	_, err = time.Parse("2006-01-02", getShortenStats.From)
	if err != nil {
		return apperror.BadRequest.WithMessage("from is invalid, expected 2006-01-02")
	}

	_, err = time.Parse("2006-01-02", getShortenStats.To)
	if err != nil {
		return apperror.BadRequest.WithMessage("to is invalid, expected 2006-01-02")
	}

	switch getShortenStats.Unit {
	case domain.UnitHour, domain.UnitDay, domain.UnitWeek, domain.UnitMonth, domain.UnitYear:
	default:
		return apperror.BadRequest.WithMessage("unit is invalid, expected (hour, day, week, month, year)")
	}

	return nil
}

func (exportShortenStats ExportShortenStats) Validate() error {
	var err error
	_, err = time.Parse("2006-01-02", exportShortenStats.From)
	if err != nil {
		return apperror.BadRequest.WithMessage("from is invalid, expected 2006-01-02")
	}

	_, err = time.Parse("2006-01-02", exportShortenStats.To)
	if err != nil {
		return apperror.BadRequest.WithMessage("to is invalid, expected 2006-01-02")
	}

	return nil
}
