package model

import (
	"cc/app/internal/domain"
	"cc/app/pkg/base62"
	"github.com/google/uuid"
	"time"
)

type Shorten struct {
	ID        uint64    `db:"id"`
	URL       string    `db:"url"`
	UserID    uuid.UUID `db:"user_id"`
	Title     string    `db:"title"`
	Tags      []string  `db:"tags"`
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
		Tags:      s.Tags,
		CreatedAt: s.CreatedAt.Unix(),
		UpdatedAt: s.UpdatedAt.Unix(),
	}
}

func (shortens Shortens) Domain(host string) domain.Shortens {
	res := make(domain.Shortens, len(shortens))

	for i, shorten := range shortens {
		res[i] = shorten.Domain(host)
	}

	return res
}
