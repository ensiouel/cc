package dto

import "time"

type CreateClick struct {
	ShortenID uint64    `json:"shorten_id"`
	Platform  string    `json:"platform"`
	OS        string    `json:"os"`
	Referrer  string    `json:"referrer"`
	IP        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}
