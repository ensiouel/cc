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

type UserStorage interface {
	CreateUser(ctx context.Context, user model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByName(ctx context.Context, name string) (model.User, error)
	ExistsUserByName(ctx context.Context, name string) (bool, error)
}

type userStorage struct {
	client postgres.Client
}

func NewUserStorage(client postgres.Client) UserStorage {
	return &userStorage{client: client}
}

func (storage *userStorage) CreateUser(ctx context.Context, user model.User) error {
	q := `
INSERT INTO 
    users (id, name, password) 
VALUES 
	($1, $2, $3)
`

	_, err := storage.client.Exec(ctx, q, user.ID, user.Name, user.Password)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *userStorage) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	q := `
SELECT 
    id, name, password 
FROM 
    users 
WHERE 
    id = $1
`

	var user model.User
	err := storage.client.Get(ctx, &user, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, apperror.NotFound
		}

		return user, apperror.Internal.WithError(err)
	}

	return user, nil
}

func (storage *userStorage) GetByName(ctx context.Context, name string) (model.User, error) {
	q := `
SELECT 
    id, name, password 
FROM 
    users 
WHERE 
    name = $1
`

	var user model.User
	err := storage.client.Get(ctx, &user, q, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, apperror.NotFound
		}

		return user, apperror.Internal.WithError(err)
	}

	return user, nil
}

func (storage *userStorage) ExistsUserByName(ctx context.Context, name string) (bool, error) {
	q := `
SELECT 
    EXISTS (
		SELECT 
			1 
		FROM 
			users 
		WHERE 
			name = $1
	)
`

	var exists bool
	err := storage.client.Get(ctx, &exists, q, name)
	if err != nil {
		return exists, apperror.Internal.WithError(err)
	}

	return exists, nil
}
