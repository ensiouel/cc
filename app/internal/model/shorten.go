package model

import (
	"cc/app/internal/domain"
	"cc/app/pkg/base62"
	"github.com/google/uuid"
	"time"
)

type Shorten struct {
	ID        uint64    `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Title     string    `db:"title"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Shortens []Shorten

func (s Shorten) Domain(host string) domain.Shorten {
	id := base62.Encode(s.ID)

	return domain.Shorten{
		ID:        id,
		Title:     s.Title,
		LongURL:   s.URL,
		ShortURL:  host + "/" + id,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (s Shortens) Domain(host string) []domain.Shorten {
	if len(s) == 0 {
		return []domain.Shorten{}
	}

	shortens := make([]domain.Shorten, len(s))

	for i, shorten := range s {
		shortens[i] = shorten.Domain(host)
	}

	return shortens
}
