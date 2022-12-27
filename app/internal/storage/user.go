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

type UserStorage interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetUserByName(ctx context.Context, name string) (model.User, error)
	ExistsUserByName(ctx context.Context, name string) (bool, error)
}

type userStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) UserStorage {
	return &userStorage{db: db}
}

func (storage *userStorage) CreateUser(ctx context.Context, user model.User) (err error) {
	q := `
INSERT INTO 
    users (id, name, password) 
VALUES 
	($1, $2, $3)
`

	_, err = storage.db.ExecContext(ctx, q, user.ID, user.Name, user.Password)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *userStorage) GetUserByID(ctx context.Context, id uuid.UUID) (user model.User, err error) {
	q := `
SELECT 
    id, name, password 
FROM 
    users 
WHERE 
    id = $1
`

	err = storage.db.GetContext(ctx, &user, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, apperror.ErrNotExists
		}

		return user, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *userStorage) GetUserByName(ctx context.Context, name string) (user model.User, err error) {
	q := `
SELECT 
    id, name, password 
FROM 
    users 
WHERE 
    name = $1
`

	err = storage.db.GetContext(ctx, &user, q, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, apperror.ErrNotExists
		}

		return user, apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *userStorage) ExistsUserByName(ctx context.Context, name string) (exists bool, err error) {
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

	err = storage.db.GetContext(ctx, &exists, q, name)
	if err != nil {
		return exists, apperror.ErrInternalError.SetError(err)
	}

	return
}
