package domain

import "time"

type Shorten struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	LongURL   string    `json:"long_url"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}
