package storage

import (
	"cc/app/internal/apperror"
	"cc/app/pkg/postgres"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type TagStorage interface {
	SelectByUser(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type tagStorage struct {
	client postgres.Client
}

func NewTagStorage(client postgres.Client) TagStorage {
	return &tagStorage{client: client}
}

func (storage *tagStorage) SelectByUser(ctx context.Context, userID uuid.UUID) (tags []string, err error) {
	q := `
SELECT DISTINCT UNNEST(tags) as title
FROM shortens
WHERE user_id = $1
`

	err = storage.client.Select(ctx, &tags, q, userID)
	if err != nil && errors.Is(err, pgx.ErrNoRows) == false {
		return tags, apperror.Internal.WithError(err)
	}

	return
}
