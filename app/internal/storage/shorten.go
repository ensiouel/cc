package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/model"
	"cc/app/pkg/postgres"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ShortenStorage interface {
	Create(ctx context.Context, shorten model.Shorten) error
	Delete(ctx context.Context, userID uuid.UUID, shortenID uint64) error

	Update(ctx context.Context, shorten model.Shorten) error

	GetByID(ctx context.Context, id uint64) (model.Shorten, error)
	GetByURL(ctx context.Context, url string) (model.Shorten, error)

	SelectByUser(ctx context.Context, userID uuid.UUID) (model.Shortens, error)
	SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (model.Shortens, error)

	GetURL(ctx context.Context, shortenID uint64) (string, error)

	ExistsByID(ctx context.Context, userID uuid.UUID, id uint64) (bool, error)
	ExistsByURL(ctx context.Context, userID uuid.UUID, url string) (bool, error)
}

type shortenStorage struct {
	client postgres.Client
}

func NewShortenStorage(client postgres.Client) ShortenStorage {
	return &shortenStorage{client: client}
}

func (storage *shortenStorage) Create(ctx context.Context, shorten model.Shorten) (err error) {
	q := `
INSERT INTO 
    shortens (id, url, user_id, title, created_at, updated_at, tags) 
VALUES 
    ($1, $2, $3, $4, $5, $6)
`

	_, err = storage.client.Exec(ctx, q,
		shorten.ID,
		shorten.URL,
		shorten.UserID,
		shorten.Title,
		shorten.CreatedAt,
		shorten.UpdatedAt,
		shorten.Tags,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) Update(ctx context.Context, shorten model.Shorten) (err error) {
	q := `
UPDATE
    shortens
SET url        = $1,
    title      = $2,
    created_at = $3,
    updated_at = $4,
    tags       = $5
WHERE id = $6
  AND user_id = $7;
`

	_, err = storage.client.Exec(ctx, q,
		shorten.URL,
		shorten.Title,
		shorten.CreatedAt,
		shorten.UpdatedAt,
		shorten.Tags,
		shorten.ID,
		shorten.UserID,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) Delete(ctx context.Context, userID uuid.UUID, shortenID uint64) (err error) {
	q := `
DELETE FROM 
	shortens
WHERE
	user_id = $1 AND 
    id = $2
`

	_, err = storage.client.Exec(ctx, q, userID, shortenID)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) GetByID(ctx context.Context, id uint64) (shorten model.Shorten, err error) {
	q := `
SELECT id,
       url,
       user_id,
       title,
       tags,
       created_at,
       updated_at
FROM shortens
WHERE id = $1
`

	err = storage.client.Get(ctx, &shorten, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shorten, apperror.NotExists.WithMessage("shorten with this id not exists")
		}

		return shorten, apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) GetByURL(ctx context.Context, url string) (shorten model.Shorten, err error) {
	q := `
SELECT id,
       url,
       user_id,
       title,
       tags,
       created_at,
       updated_at
FROM shortens
WHERE url = $1
`

	err = storage.client.Get(ctx, &shorten, q, url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shorten, apperror.NotExists.WithMessage("shorten with this url not exists")
		}

		return shorten, apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) GetURL(ctx context.Context, shortenID uint64) (url string, err error) {
	q := `
SELECT 
    url 
FROM 
    shortens 
WHERE 
    id = $1
`

	err = storage.client.Get(ctx, &url, q, shortenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return url, apperror.NotExists.WithMessage("url with this shorten_id not exists")
		}

		return url, apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) SelectByUser(ctx context.Context, id uuid.UUID) (shortens model.Shortens, err error) {
	q := `
SELECT id,
       url,
       user_id,
       title,
       tags,
       created_at,
       updated_at
FROM shortens
WHERE user_id = $1
`

	err = storage.client.Select(ctx, &shortens, q, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) == false {
		return shortens, apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (shortens model.Shortens, err error) {
	q := `
SELECT id,
       url,
       user_id,
       title,
       tags,
       created_at,
       updated_at
FROM shortens
WHERE user_id = $1
  AND tags @> $2
`

	err = storage.client.Select(ctx, &shortens, q, userID, tags)
	if err != nil && errors.Is(err, pgx.ErrNoRows) == false {
		return shortens, apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) ExistsByID(ctx context.Context, userID uuid.UUID, id uint64) (exists bool, err error) {
	q := `
SELECT 
    EXISTS (
		SELECT 
			1 
		FROM 
			shortens 
		WHERE 
			id = $1 AND 
			user_id = $2
	)
`

	err = storage.client.Get(ctx, &exists, q, id, userID)
	if err != nil {
		return exists, apperror.Internal.WithError(err)
	}

	return
}

func (storage *shortenStorage) ExistsByURL(ctx context.Context, userID uuid.UUID, url string) (exists bool, err error) {
	q := `
SELECT 
    EXISTS (
		SELECT 
			1 
		FROM 
			shortens 
		WHERE 
			url = $1 AND 
			user_id = $2
	)
`

	err = storage.client.Get(ctx, &exists, q, url, userID)
	if err != nil {
		return exists, apperror.Internal.WithError(err)
	}

	return
}
