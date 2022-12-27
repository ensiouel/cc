package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShortenStorage interface {
	CreateShorten(ctx context.Context, shorten model.Shorten) error
	UpdateShorten(ctx context.Context, shorten model.Shorten) error
	DeleteShorten(ctx context.Context, userID uuid.UUID, shortenID uint64) error
	GetShortenByID(ctx context.Context, id uint64) (model.Shorten, error)
	GetShortenByURL(ctx context.Context, url string) (model.Shorten, error)
	SelectShortensByUserID(ctx context.Context, id uuid.UUID) ([]model.Shorten, error)
	ExistsShortenByID(ctx context.Context, userID uuid.UUID, id uint64) (bool, error)
	ExistsShortenByURL(ctx context.Context, userID uuid.UUID, url string) (bool, error)
}

type shortenStorage struct {
	db *sqlx.DB
}

func NewShortenStorage(db *sqlx.DB) ShortenStorage {
	return &shortenStorage{db: db}
}

func (storage *shortenStorage) CreateShorten(ctx context.Context, shorten model.Shorten) (err error) {
	q := `
INSERT INTO 
    shortens (id, url, user_id, title, created_at) 
VALUES 
    ($1, $2, $3, $4, $5)
`

	_, err = storage.db.ExecContext(ctx, q,
		shorten.ID,
		shorten.URL,
		shorten.UserID,
		shorten.Title,
		shorten.CreatedAt,
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
    created_at = $3
WHERE
    id = $4 AND
    user_id = $5
`

	_, err = storage.db.ExecContext(ctx, q,
		shorten.URL,
		shorten.Title,
		shorten.CreatedAt,
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

	_, err = storage.db.ExecContext(ctx, q, userID, shortenID)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) GetShortenByID(ctx context.Context, id uint64) (shorten model.Shorten, err error) {
	q := `
SELECT 
    id, url, user_id, title, created_at 
FROM 
    shortens 
WHERE
	id = $1
`

	err = storage.db.GetContext(ctx, &shorten, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shorten, apperror.ErrNotExists
		}

		return shorten, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) GetShortenByURL(ctx context.Context, url string) (shorten model.Shorten, err error) {
	q := `
SELECT 
    id, url, user_id, title, created_at 
FROM 
    shortens 
WHERE
	url = $1
`

	err = storage.db.GetContext(ctx, &shorten, q, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shorten, apperror.ErrNotExists
		}

		return shorten, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *shortenStorage) SelectShortensByUserID(ctx context.Context, id uuid.UUID) (shortens []model.Shorten, err error) {
	q := `
SELECT 
    id, url, user_id, title, created_at 
FROM
    shortens
WHERE
    user_id = $1
`

	err = storage.db.SelectContext(ctx, &shortens, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return shortens, apperror.ErrNotExists
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

	err = storage.db.GetContext(ctx, &exists, q, id, userID)
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

	err = storage.db.GetContext(ctx, &exists, q, url, userID)
	if err != nil {
		return exists, apperror.ErrInternalError.SetError(err)
	}

	return
}
