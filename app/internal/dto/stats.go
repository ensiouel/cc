package dto

import "time"

type CreateClick struct {
	ShortenID uint64    `json:"shorten_id"`
	Platform  string    `json:"platform"`
	OS        string    `json:"os"`
	Referer   string    `json:"referer"`
	IP        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}
