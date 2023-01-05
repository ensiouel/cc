package domain

import (
	"time"
)

const (
	UnitMinute Unit = "minute"
	UnitHour   Unit = "hour"
	UnitDay    Unit = "day"
	UnitWeek   Unit = "week"
	UnitMonth  Unit = "month"
	UnitYear   Unit = "year"
)

type Unit string

type Click struct {
	ShortenID uint64    `json:"shorten_id"`
	Platform  string    `json:"platform"`
	OS        string    `json:"os"`
	Referrer  string    `json:"referrer"`
	Timestamp time.Time `json:"timestamp"`
}

type ClickState struct {
	Count int       `json:"count"`
	Date  time.Time `json:"date"`
}

type MetricState struct {
	Name  string    `json:"name"`
	Count int       `json:"count"`
	Date  time.Time `json:"date"`
}

type MetricSummaryState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ClickStats struct {
	Clicks []ClickState `json:"clicks"`
	Unit   Unit         `json:"unit"`
	Units  int          `json:"units"`
}

type ClickSummaryStats struct {
	Clicks int `json:"clicks"`
}

type MetricStats struct {
	Metrics []MetricState `json:"metrics"`
	Unit    Unit          `json:"unit"`
	Units   int           `json:"units"`
}

type MetricSummaryStats struct {
	Metrics []MetricSummaryState `json:"metrics"`
}
