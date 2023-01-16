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
	CreateShorten(ctx context.Context, shorten model.Shorten) error
	UpdateShorten(ctx context.Context, shorten model.Shorten) error
	DeleteShorten(ctx context.Context, userID uuid.UUID, shortenID uint64) error
	GetShortenByID(ctx context.Context, id uint64) (model.Shorten, error)
	GetShortenByURL(ctx context.Context, url string) (model.Shorten, error)
	GetShortenURL(ctx context.Context, shortenID uint64) (string, error)
	SelectShortensByUserID(ctx context.Context, id uuid.UUID) ([]model.Shorten, error)
	ExistsShortenByID(ctx context.Context, userID uuid.UUID, id uint64) (bool, error)
	ExistsShortenByURL(ctx context.Context, userID uuid.UUID, url string) (bool, error)
}

type shortenStorage struct {
	client postgres.Client
}

func NewShortenStorage(client postgres.Client) ShortenStorage {
	return &shortenStorage{client: client}
}

func (storage *shortenStorage) CreateShorten(ctx context.Context, shorten model.Shorten) (err error) {
	q := `
INSERT INTO 
    shortens (id, url, user_id, title, created_at, updated_at) 
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
	)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) UpdateShorten(ctx context.Context, shorten model.Shorten) (err error) {
	q := `
UPDATE 
    shortens
SET 
    url = $1,
    title = $2,
    created_at = $3,
    updated_at = $4
WHERE
    id = $5 AND
    user_id = $6
`

	_, err = storage.client.Exec(ctx, q,
		shorten.URL,
		shorten.Title,
		shorten.CreatedAt,
		shorten.UpdatedAt,
		shorten.ID,
		shorten.UserID,
	)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) DeleteShorten(ctx context.Context, userID uuid.UUID, shortenID uint64) (err error) {
	q := `
DELETE FROM 
	shortens
WHERE
	user_id = $1 AND 
    id = $2
`

	_, err = storage.client.Exec(ctx, q, userID, shortenID)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) GetShortenByID(ctx context.Context, id uint64) (shorten model.Shorten, err error) {
	q := `
SELECT 
    id, url, user_id, title, created_at, updated_at 
FROM 
    shortens 
WHERE
	id = $1
`

	err = storage.client.Get(ctx, &shorten, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shorten, apperror.ErrNotExists.SetMessage("shorten with this id not exists")
		}

		return shorten, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) GetShortenByURL(ctx context.Context, url string) (shorten model.Shorten, err error) {
	q := `
SELECT 
    id, url, user_id, title, created_at, updated_at
FROM 
    shortens 
WHERE
	url = $1
`

	err = storage.client.Get(ctx, &shorten, q, url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shorten, apperror.ErrNotExists.SetMessage("shorten with this url not exists")
		}

		return shorten, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) GetShortenURL(ctx context.Context, shortenID uint64) (url string, err error) {
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
			return url, apperror.ErrNotExists.SetMessage("url with this shorten_id not exists")
		}

		return url, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) SelectShortensByUserID(ctx context.Context, id uuid.UUID) (shortens []model.Shorten, err error) {
	q := `
SELECT 
    id, url, user_id, title, created_at, updated_at
FROM
    shortens
WHERE
    user_id = $1
`

	err = storage.client.Select(ctx, &shortens, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shortens, apperror.ErrNotExists.SetMessage("shorten with this user_id not exists")
		}

		return shortens, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) ExistsShortenByID(ctx context.Context, userID uuid.UUID, id uint64) (exists bool, err error) {
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
		return exists, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) ExistsShortenByURL(ctx context.Context, userID uuid.UUID, url string) (exists bool, err error) {
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
		return exists, apperror.ErrInternalError.SetError(err)
	}

	return
}
