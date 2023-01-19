package domain

import (
	"time"
)

const (
	UnitHour  Unit = "hour"
	UnitDay   Unit = "day"
	UnitWeek  Unit = "week"
	UnitMonth Unit = "month"
	UnitYear  Unit = "year"
)

type Unit string

type Click struct {
	ShortenID uint64    `json:"shorten_id"`
	Platform  string    `json:"platform"`
	OS        string    `json:"os"`
	Referrer  string    `json:"referrer"`
	Timestamp time.Time `json:"timestamp"`
}

type Stats struct {
	Click    ClickMetric `json:"click"`
	Platform []Metric    `json:"platform"`
	OS       []Metric    `json:"os"`
	Referer  []Metric    `json:"referer"`
}

type ClickMetric struct {
	Total  int `json:"total"`
	Values []struct {
		Timestamp time.Time `json:"timestamp"`
		Count     int       `json:"count"`
	} `json:"values"`
}

type Metric struct {
	Name   string `json:"name"`
	Total  int    `json:"total"`
	Values []struct {
		Timestamp time.Time `json:"timestamp"`
		Count     int       `json:"count"`
	} `json:"values"`
}
