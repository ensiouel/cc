package model

import (
	"cc/internal/domain"
	"encoding/json"
	"time"
)

type Click struct {
	ShortenID uint64    `db:"shorten_id"`
	Platform  string    `db:"platform"`
	OS        string    `db:"os"`
	Referer   string    `db:"referer"`
	IP        string    `db:"ip"`
	Timestamp time.Time `db:"timestamp"`
}

type Clicks []Click

type ClickMetric struct {
	Total  int `db:"total"`
	Diff   int `db:"diff"`
	Values []byte
}

type Metric struct {
	Name   string `db:"name"`
	Total  int    `db:"total"`
	Diff   int    `db:"diff"`
	Values []byte
}

type Metrics []Metric

func (m ClickMetric) Domain() (metric domain.ClickMetric) {
	metric.Total = m.Total
	metric.Diff = m.Diff

	_ = json.Unmarshal(m.Values, &metric.Values)

	return
}

func (m Metric) Domain() (metric domain.Metric) {
	metric.Name = m.Name
	metric.Total = m.Total
	metric.Diff = m.Diff

	_ = json.Unmarshal(m.Values, &metric.Values)

	return
}

func (m Metrics) Domain() []domain.Metric {
	metrics := make([]domain.Metric, len(m))

	for i, v := range m {
		metrics[i] = v.Domain()
	}

	return metrics
}

func (c Click) Domain() domain.Click {
	return domain.Click{
		ShortenID: c.ShortenID,
		Platform:  c.Platform,
		OS:        c.OS,
		Referer:   c.Referer,
		Timestamp: c.Timestamp,
	}
}

func (m Clicks) Domain() []domain.Click {
	metrics := make([]domain.Click, len(m))

	for i, v := range m {
		metrics[i] = v.Domain()
	}

	return metrics
}
