package model

import (
	"cc/app/internal/domain"
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

type ClickState struct {
	Count int       `db:"count"`
	Date  time.Time `db:"date"`
}

type MetricState struct {
	Name  string    `db:"name"`
	Count int       `db:"count"`
	Date  time.Time `db:"date"`
}

type MetricSummaryState struct {
	Name  string `db:"name"`
	Count int    `db:"count"`
}

type ClickStats []ClickState

type MetricStats []MetricState

type MetricSummaryStats []MetricSummaryState

func (clickState ClickState) Domain() domain.ClickState {
	return domain.ClickState{
		Count: clickState.Count,
		Date:  clickState.Date,
	}
}

func (metricState MetricState) Domain() domain.MetricState {
	return domain.MetricState{
		Name:  metricState.Name,
		Count: metricState.Count,
		Date:  metricState.Date,
	}
}

func (metricSummaryState MetricSummaryState) Domain() domain.SummaryMetricState {
	return domain.SummaryMetricState{
		Name:  metricSummaryState.Name,
		Count: metricSummaryState.Count,
	}
}

func (s ClickStats) Domain() []domain.ClickState {
	if len(s) == 0 {
		return []domain.ClickState{}
	}

	clickStats := make([]domain.ClickState, len(s))

	for i, shorten := range s {
		clickStats[i] = shorten.Domain()
	}

	return clickStats
}

func (s MetricStats) Domain() []domain.MetricState {
	if len(s) == 0 {
		return []domain.MetricState{}
	}

	clickStats := make([]domain.MetricState, len(s))

	for i, shorten := range s {
		clickStats[i] = shorten.Domain()
	}

	return clickStats
}

func (s MetricSummaryStats) Domain() []domain.SummaryMetricState {
	if len(s) == 0 {
		return []domain.SummaryMetricState{}
	}

	clickStats := make([]domain.SummaryMetricState, len(s))

	for i, shorten := range s {
		clickStats[i] = shorten.Domain()
	}

	return clickStats
}
