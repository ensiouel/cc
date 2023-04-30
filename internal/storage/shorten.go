package storage

import (
	"cc/internal/model"
	"cc/pkg/apperror"
	"cc/pkg/postgres"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate moq -out shorten_mock.go . ShortenStorage
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

func (storage *shortenStorage) Create(ctx context.Context, shorten model.Shorten) error {
	q := `
INSERT INTO 
    shortens (id, url, user_id, title, created_at, updated_at, tags) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7)
`

	_, err := storage.client.Exec(ctx, q,
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

	return nil
}

func (storage *shortenStorage) Update(ctx context.Context, shorten model.Shorten) error {
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

	_, err := storage.client.Exec(ctx, q,
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

	return nil
}

func (storage *shortenStorage) Delete(ctx context.Context, userID uuid.UUID, shortenID uint64) error {
	q := `
DELETE FROM 
	shortens
WHERE
	user_id = $1 AND 
    id = $2
`

	_, err := storage.client.Exec(ctx, q, userID, shortenID)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *shortenStorage) GetByID(ctx context.Context, id uint64) (model.Shorten, error) {
	return storage.getBy(ctx, "id", id)
}

func (storage *shortenStorage) GetByURL(ctx context.Context, url string) (model.Shorten, error) {
	return storage.getBy(ctx, "url", url)
}

func (storage *shortenStorage) GetURL(ctx context.Context, shortenID uint64) (string, error) {
	q := `
SELECT 
    url 
FROM 
    shortens 
WHERE 
    id = $1
`

	var url string
	err := storage.client.Get(ctx, &url, q, shortenID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return url, apperror.NotFound
		}

		return url, apperror.Internal.WithError(err)
	}

	return url, nil
}

func (storage *shortenStorage) SelectByUser(ctx context.Context, id uuid.UUID) (model.Shortens, error) {
	return storage.selectBy(ctx, "user_id", id)
}

func (storage *shortenStorage) SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (model.Shortens, error) {
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

	var shortens model.Shortens
	err := storage.client.Select(ctx, &shortens, q, userID, tags)
	if err != nil && errors.Is(err, pgx.ErrNoRows) == false {
		return shortens, apperror.Internal.WithError(err)
	}

	return shortens, nil
}

func (storage *shortenStorage) ExistsByID(ctx context.Context, userID uuid.UUID, id uint64) (bool, error) {
	return storage.existsBy(ctx, userID, "id", id)
}

func (storage *shortenStorage) ExistsByURL(ctx context.Context, userID uuid.UUID, url string) (bool, error) {
	return storage.existsBy(ctx, userID, "url", url)
}

func (storage *shortenStorage) getBy(ctx context.Context, column string, value any) (model.Shorten, error) {
	q := `
SELECT id,
       url,
       user_id,
       title,
       tags,
       created_at,
       updated_at
FROM shortens
WHERE ` + column + ` = $1`

	var shorten model.Shorten
	err := storage.client.Get(ctx, &shorten, q, value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shorten, apperror.NotFound.WithError(err)
		}

		return shorten, apperror.Internal.WithError(err)
	}

	return shorten, nil
}

func (storage *shortenStorage) selectBy(ctx context.Context, column string, value any) (model.Shortens, error) {
	q := `
SELECT id,
       url,
       user_id,
       title,
       tags,
       created_at,
       updated_at
FROM shortens
WHERE ` + column + ` = $1`

	var shortens model.Shortens
	err := storage.client.Select(ctx, &shortens, q, value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shortens, apperror.NotFound.WithError(err)
		}

		return shortens, apperror.Internal.WithError(err)
	}

	return shortens, nil
}

func (storage *shortenStorage) existsBy(ctx context.Context, userID uuid.UUID, column string, value any) (bool, error) {
	q := `
SELECT 
    EXISTS (
		SELECT 
			1 
		FROM 
			shortens 
		WHERE 
			` + column + ` = $1 AND 
			user_id = $2
	)
`

	var exists bool
	err := storage.client.Get(ctx, &exists, q, value, userID)
	if err != nil {
		return exists, apperror.Internal.WithError(err)
	}

	return exists, nil
}
