package model

import (
	"cc/app/internal/domain"
	"encoding/json"
	"time"
)

type Click struct {
	ShortenID uint64    `db:"shorten_id"`
	Platform  string    `db:"platform"`
	OS        string    `db:"os"`
	Referrer  string    `db:"referrer"`
	IP        string    `db:"ip"`
	Timestamp time.Time `db:"timestamp"`
}

type ClickMetric struct {
	Total  int `db:"total"`
	Values []byte
}

type Metric struct {
	Name   string `db:"name"`
	Total  int    `db:"total"`
	Values []byte
}

type Metrics []Metric

func (m ClickMetric) Domain() (clickMetric domain.ClickMetric) {
	clickMetric.Total = m.Total

	_ = json.Unmarshal(m.Values, &clickMetric.Values)

	return
}

func (m Metric) Domain() (metric domain.Metric) {
	metric.Name = m.Name
	metric.Total = m.Total

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
